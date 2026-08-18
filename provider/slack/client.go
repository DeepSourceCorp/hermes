package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"

	log "github.com/sirupsen/logrus"
)

const postMessageURL = "https://slack.com/api/chat.postMessage"

type Client struct {
	HTTPClient provider.IHTTPClient

	// Sleep is swapped out in tests so backoff does not slow them down.
	Sleep func(time.Duration)
}

// sleep waits for d, or gives up early if the caller has gone away or the
// listing budget has run out. It reports the context error in that case.
func (c *Client) sleep(ctx context.Context, d time.Duration) error {
	if c.Sleep != nil {
		c.Sleep(d)
		return ctx.Err()
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-timer.C:
		return ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
}

type SendMessageRequest struct {
	Channel     string      `json:"channel"`
	Blocks      interface{} `json:"blocks,omitempty"`
	Text        string      `json:"text,omitempty"`
	BearerToken string      `json:"-"`
}

type SendMessageResponse struct {
	Ok bool `json:"ok"`
}

func (c *Client) SendMessage(request *SendMessageRequest) (*SendMessageResponse, domain.IError) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		log.Errorf("slack: failed encoding request: %v", err)
		return nil, errFailedSendPermanent(err.Error())
	}

	req, err := http.NewRequest("POST", postMessageURL, &buf)
	if err != nil {
		log.Errorf("slack: sending request: %v", err)
		return nil, errFailedSendPermanent(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("slack: something went wrong sending request: %v", err)
		return nil, errFailedSendTemporary("something went wrong while sending messsage to slack")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(SendMessageResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Errorf("slack: failed decoding response: %v", err)
		return nil, errFailedSendPermanent(err.Error())
	}
	return response, nil
}

func handleHTTPFailure(response *http.Response) domain.IError {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("slack: failed reading response body: %v", err)
		return errFailedSendPermanent(err.Error())
	}

	body := strings.TrimSpace(string(b))

	// Rate limits and server errors are worth retrying. Everything else is a
	// permanent failure for this request.
	if response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500 {
		log.Errorf("slack: retryable failure, status=%d error=%s", response.StatusCode, body)
		return errFailedSendTemporary(fmt.Sprintf("received %d, error=%s", response.StatusCode, body))
	}

	log.Errorf("slack: permanent failure, status=%d error=%s", response.StatusCode, body)
	return errFailedSendPermanent(fmt.Sprintf("received %d, error=%s", response.StatusCode, body))
}

const getChannelsURL = "https://slack.com/api/conversations.list?types=public_channel,private_channel&exclude_archived=true&limit=1000"

const (
	// conversations.list is a Slack Tier 2 method, which allows roughly 20
	// requests a minute. A workspace large enough to need several pages will
	// trip that limit mid-pagination, so the listing has to back off and retry
	// instead of abandoning the whole thing.
	maxRateLimitRetries = 2

	// Slack's Retry-After for Tier 2 methods is often 30s or more. This runs
	// inside the synchronous OAuth callback, so the wait is capped and partial
	// results are preferred over holding the request open indefinitely.
	maxRetryAfter     = 15 * time.Second
	defaultRetryAfter = 5 * time.Second

	// channelListingBudget bounds the listing as a whole. Retrying per page is
	// not enough on a workspace big enough to stay rate limited across dozens of
	// pages: the listing then runs for minutes, long after the caller has timed
	// out and stopped waiting for it, while still spending the workspace's Slack
	// quota and starving the retry the user just kicked off. The budget has to
	// stay comfortably under the caller's own timeout so a partial list arrives
	// before it gives up.
	channelListingBudget = 20 * time.Second

	// Runaway guard in case Slack keeps handing back a next_cursor.
	maxChannelPages = 200

	// The `error` Slack sets on a rate limited response body.
	slackRateLimitedError = "ratelimited"
)

type GetChannelsRequest struct {
	BearerToken string `json:"_"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type GetChannelsResponse struct {
	Ok               bool             `json:"ok"`
	Error            string           `json:"error"`
	Channels         []Channel        `json:"channels"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
}

// rateLimitedError marks a channel page fetch that Slack rate limited, and
// carries the wait Slack asked for so pagination can retry the same cursor.
type rateLimitedError struct {
	domain.IError
	retryAfter time.Duration
}

func newRateLimitedError(retryAfter time.Duration, internal string) *rateLimitedError {
	return &rateLimitedError{IError: errFailedOptsFetch(internal), retryAfter: retryAfter}
}

func isRateLimited(err domain.IError) bool {
	_, ok := err.(*rateLimitedError)
	return ok
}

// abandonedError marks a listing that stopped because the budget ran out or
// the caller disconnected, rather than because Slack refused us.
type abandonedError struct {
	domain.IError
}

func newAbandonedError(internal string) *abandonedError {
	return &abandonedError{IError: errFailedOptsFetch(internal)}
}

func isAbandoned(err domain.IError) bool {
	_, ok := err.(*abandonedError)
	return ok
}

// partialResultsUsable reports whether the pages fetched before err are still
// worth returning. Slack rate limiting us and the budget running out both stop
// the listing short without invalidating what already came back; anything else
// means we cannot trust the listing at all.
func partialResultsUsable(err domain.IError) bool {
	return isRateLimited(err) || isAbandoned(err)
}

// retryAfterFrom reads Slack's Retry-After header, falling back to a default
// when it is missing or unparseable, and clamping it to maxRetryAfter.
func retryAfterFrom(header http.Header) time.Duration {
	seconds, err := strconv.Atoi(strings.TrimSpace(header.Get("Retry-After")))
	if err != nil || seconds <= 0 {
		return defaultRetryAfter
	}

	if retryAfter := time.Duration(seconds) * time.Second; retryAfter < maxRetryAfter {
		return retryAfter
	}
	return maxRetryAfter
}

func (c *Client) getChannelsPage(ctx context.Context, request *GetChannelsRequest, cursor string) (*GetChannelsResponse, domain.IError) {
	var response = new(GetChannelsResponse)

	requestUrl := getChannelsURL
	if cursor != "" {
		requestUrl += "&cursor=" + url.QueryEscape(cursor)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", requestUrl, http.NoBody)
	if err != nil {
		log.Errorf("slack: failed creating request for options: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		// A cancelled request is the budget or the caller, not Slack, so it must
		// not be logged or classified as a fetch failure.
		if ctxErr := ctx.Err(); ctxErr != nil {
			return response, newAbandonedError(ctxErr.Error())
		}
		log.Errorf("slack: failed sending request for options: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		b, _ := io.ReadAll(resp.Body)
		return response, newRateLimitedError(
			retryAfterFrom(resp.Header),
			fmt.Sprintf("slack rate limited the channel listing, error=%s", strings.TrimSpace(string(b))),
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		log.Errorf("slack: non-2xx response while fetching options: status=%d", resp.StatusCode)
		return response, handleHTTPFailure(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Errorf("slack: failed decoding options response: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}

	// Slack reports application level failures as `ok: false` on a 200, so the
	// status code alone is not enough to tell whether the page came back.
	if !response.Ok {
		if response.Error == slackRateLimitedError {
			return response, newRateLimitedError(
				retryAfterFrom(resp.Header),
				"slack rate limited the channel listing",
			)
		}
		log.Errorf("slack: channel listing failed with error=%s", response.Error)
		return response, errFailedOptsFetch(fmt.Sprintf("slack returned error=%s", response.Error))
	}

	return response, nil
}

// getChannelsPageWithBackoff fetches a single page, retrying a bounded number
// of times while Slack rate limits us. Non rate limit failures are returned
// straight away.
func (c *Client) getChannelsPageWithBackoff(ctx context.Context, request *GetChannelsRequest, cursor string) (*GetChannelsResponse, domain.IError) {
	for attempt := 0; ; attempt++ {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return new(GetChannelsResponse), newAbandonedError(ctxErr.Error())
		}

		response, err := c.getChannelsPage(ctx, request, cursor)
		if err == nil {
			return response, nil
		}

		rateLimited, ok := err.(*rateLimitedError)
		if !ok {
			return response, err
		}

		if attempt >= maxRateLimitRetries {
			return response, rateLimited
		}

		log.Warnf(
			"slack: rate limited fetching channel page %q, retrying in %v (attempt %d of %d)",
			cursor, rateLimited.retryAfter, attempt+1, maxRateLimitRetries,
		)

		// Waiting out the whole Retry-After is pointless once the budget is
		// spent, since nothing would be left to fetch the page with.
		if sleepErr := c.sleep(ctx, rateLimited.retryAfter); sleepErr != nil {
			return response, newAbandonedError(sleepErr.Error())
		}
	}
}

func (c *Client) GetChannels(ctx context.Context, request *GetChannelsRequest) ([]map[string]string, domain.IError) {
	ctx, cancel := context.WithTimeout(ctx, channelListingBudget)
	defer cancel()

	channels := make([]map[string]string, 0)
	cursor := ""

	for page := 0; page < maxChannelPages; page++ {
		response, err := c.getChannelsPageWithBackoff(ctx, request, cursor)
		if err != nil {
			// Slack kept rate limiting us, or we ran out of budget. The pages
			// that did come back are far more useful than a hard failure, which
			// stops the integration from being installed at all.
			if partialResultsUsable(err) && len(channels) > 0 {
				log.Warnf(
					"slack: stopped paginating channels (%v), returning the %d channels fetched so far",
					err.Error(), len(channels),
				)
				return channels, nil
			}

			log.Errorf("slack: Error fetching page %v: %v", cursor, err)
			return channels, err
		}

		for _, v := range response.Channels {
			channels = append(channels, map[string]string{
				"id":   v.ID,
				"name": v.Name,
			})
		}

		cursor = response.ResponseMetadata.NextCursor
		if cursor == "" {
			return channels, nil
		}
	}

	log.Warnf(
		"slack: hit the %d page cap while paginating channels, returning %d channels",
		maxChannelPages, len(channels),
	)
	return channels, nil
}
