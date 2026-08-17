package slack

import (
	"bytes"
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

func (c *Client) sleep(d time.Duration) {
	if c.Sleep != nil {
		c.Sleep(d)
		return
	}
	time.Sleep(d)
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

func (c *Client) getChannelsPage(request *GetChannelsRequest, cursor string) (*GetChannelsResponse, domain.IError) {
	var response = new(GetChannelsResponse)

	requestUrl := getChannelsURL
	if cursor != "" {
		requestUrl += "&cursor=" + url.QueryEscape(cursor)
	}

	req, err := http.NewRequest("GET", requestUrl, http.NoBody)
	if err != nil {
		log.Errorf("slack: failed creating request for options: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
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
func (c *Client) getChannelsPageWithBackoff(request *GetChannelsRequest, cursor string) (*GetChannelsResponse, domain.IError) {
	for attempt := 0; ; attempt++ {
		response, err := c.getChannelsPage(request, cursor)
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
		c.sleep(rateLimited.retryAfter)
	}
}

func (c *Client) GetChannels(request *GetChannelsRequest) ([]map[string]string, domain.IError) {
	channels := make([]map[string]string, 0)
	cursor := ""

	for page := 0; page < maxChannelPages; page++ {
		response, err := c.getChannelsPageWithBackoff(request, cursor)
		if err != nil {
			// Slack kept rate limiting us. The pages that did come back are far
			// more useful than a hard failure, which stops the integration from
			// being installed at all.
			if isRateLimited(err) && len(channels) > 0 {
				log.Warnf(
					"slack: rate limited while paginating channels, returning the %d channels fetched so far",
					len(channels),
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
