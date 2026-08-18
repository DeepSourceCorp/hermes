package slack

import (
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

// stubHTTPClient replays a canned list of responses, one per request, and
// records the request URLs it was asked for.
type stubHTTPClient struct {
	responses []*http.Response
	requested []string
	calls     int

	// onCall runs after each request, so a test can cancel the listing at a
	// chosen point in the pagination.
	onCall func(calls int)
}

func (s *stubHTTPClient) Do(req *http.Request) (*http.Response, error) {
	s.requested = append(s.requested, req.URL.String())

	if s.calls >= len(s.responses) {
		s.calls++
		return nil, io.ErrUnexpectedEOF
	}

	resp := s.responses[s.calls]
	s.calls++
	if s.onCall != nil {
		s.onCall(s.calls)
	}
	return resp, nil
}

func response(status int, body string, header http.Header) *http.Response {
	if header == nil {
		header = http.Header{}
	}
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     header,
	}
}

func rateLimitedResponse(retryAfter string) *http.Response {
	header := http.Header{}
	if retryAfter != "" {
		header.Set("Retry-After", retryAfter)
	}
	return response(http.StatusTooManyRequests, `{"ok":false,"error":"ratelimited"}`, header)
}

// newTestClient returns a client whose backoff is recorded rather than slept.
func newTestClient(responses ...*http.Response) (*Client, *stubHTTPClient, *[]time.Duration) {
	stub := &stubHTTPClient{responses: responses}
	var slept []time.Duration
	client := &Client{
		HTTPClient: stub,
		Sleep:      func(d time.Duration) { slept = append(slept, d) },
	}
	return client, stub, &slept
}

func channelNames(channels []map[string]string) []string {
	names := make([]string, 0, len(channels))
	for _, channel := range channels {
		names = append(names, channel["name"])
	}
	return names
}

func TestClient_GetChannels_PaginatesUntilCursorIsEmpty(t *testing.T) {
	client, stub, _ := newTestClient(
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}],"response_metadata":{"next_cursor":"page2"}}`, nil),
		response(200, `{"ok":true,"channels":[{"id":"C2","name":"random"}]}`, nil),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v", err)
	}

	if want := []string{"general", "random"}; !reflect.DeepEqual(channelNames(got), want) {
		t.Errorf("GetChannels() = %v, want %v", channelNames(got), want)
	}
	if stub.calls != 2 {
		t.Errorf("GetChannels() made %d requests, want 2", stub.calls)
	}
	if !strings.Contains(stub.requested[1], "cursor=page2") {
		t.Errorf("second request = %q, want it to carry cursor=page2", stub.requested[1])
	}
}

// The regression this fixes: Slack rate limits a page mid-pagination and the
// whole channel listing used to be abandoned, which silently blocked the
// integration from installing.
func TestClient_GetChannels_RetriesRateLimitedPage(t *testing.T) {
	client, stub, slept := newTestClient(
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}],"response_metadata":{"next_cursor":"page2"}}`, nil),
		rateLimitedResponse("3"),
		response(200, `{"ok":true,"channels":[{"id":"C2","name":"random"}]}`, nil),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v", err)
	}

	if want := []string{"general", "random"}; !reflect.DeepEqual(channelNames(got), want) {
		t.Errorf("GetChannels() = %v, want %v", channelNames(got), want)
	}
	if want := []time.Duration{3 * time.Second}; !reflect.DeepEqual(*slept, want) {
		t.Errorf("backoff = %v, want %v", *slept, want)
	}
	// The retry must re-request the same cursor, not skip the page.
	if !strings.Contains(stub.requested[2], "cursor=page2") {
		t.Errorf("retry request = %q, want it to carry cursor=page2", stub.requested[2])
	}
}

func TestClient_GetChannels_RetriesRateLimitedBodyOn200(t *testing.T) {
	client, _, slept := newTestClient(
		response(200, `{"ok":false,"error":"ratelimited"}`, nil),
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}]}`, nil),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v", err)
	}

	if want := []string{"general"}; !reflect.DeepEqual(channelNames(got), want) {
		t.Errorf("GetChannels() = %v, want %v", channelNames(got), want)
	}
	if want := []time.Duration{defaultRetryAfter}; !reflect.DeepEqual(*slept, want) {
		t.Errorf("backoff = %v, want %v", *slept, want)
	}
}

// Once retries are exhausted, an install is still possible with the pages that
// did come back, so partial results beat a hard failure.
func TestClient_GetChannels_ReturnsPartialResultsWhenRetriesExhausted(t *testing.T) {
	client, _, slept := newTestClient(
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}],"response_metadata":{"next_cursor":"page2"}}`, nil),
		rateLimitedResponse("1"),
		rateLimitedResponse("1"),
		rateLimitedResponse("1"),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v, want partial success", err)
	}

	if want := []string{"general"}; !reflect.DeepEqual(channelNames(got), want) {
		t.Errorf("GetChannels() = %v, want %v", channelNames(got), want)
	}
	if len(*slept) != maxRateLimitRetries {
		t.Errorf("backoff attempts = %d, want %d", len(*slept), maxRateLimitRetries)
	}
}

// With nothing at all to show, the rate limit is a real failure.
func TestClient_GetChannels_ErrorsWhenRateLimitedWithNoResults(t *testing.T) {
	client, _, _ := newTestClient(
		rateLimitedResponse("1"),
		rateLimitedResponse("1"),
		rateLimitedResponse("1"),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err == nil {
		t.Fatalf("GetChannels() error = nil, want a rate limit error")
	}
	if !isRateLimited(err) {
		t.Errorf("GetChannels() error = %v, want it to be a rate limit error", err)
	}
	if len(got) != 0 {
		t.Errorf("GetChannels() = %v, want no channels", got)
	}
}

func TestClient_GetChannels_ErrorsOnSlackApplicationError(t *testing.T) {
	client, stub, _ := newTestClient(
		response(200, `{"ok":false,"error":"invalid_auth"}`, nil),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err == nil {
		t.Fatalf("GetChannels() error = nil, want an error for ok:false")
	}
	if isRateLimited(err) {
		t.Errorf("GetChannels() error = %v, want a non rate limit error", err)
	}
	if !strings.Contains(err.Error(), "invalid_auth") {
		t.Errorf("GetChannels() error = %q, want it to mention invalid_auth", err.Error())
	}
	if len(got) != 0 {
		t.Errorf("GetChannels() = %v, want no channels", got)
	}
	// An application error is not retried.
	if stub.calls != 1 {
		t.Errorf("GetChannels() made %d requests, want 1", stub.calls)
	}
}

// An empty workspace must come back as an empty list rather than a nil one, so
// callers iterating the options do not trip over a null.
func TestClient_GetChannels_ReturnsEmptySliceForNoChannels(t *testing.T) {
	client, _, _ := newTestClient(
		response(200, `{"ok":true,"channels":[]}`, nil),
	)

	got, err := client.GetChannels(context.Background(), &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v", err)
	}
	if got == nil {
		t.Fatal("GetChannels() = nil, want an empty slice")
	}
	if len(got) != 0 {
		t.Errorf("GetChannels() = %v, want no channels", got)
	}
}

func TestRetryAfterFrom(t *testing.T) {
	tests := []struct {
		name       string
		retryAfter string
		want       time.Duration
	}{
		{name: "honours the header", retryAfter: "7", want: 7 * time.Second},
		{name: "falls back when absent", retryAfter: "", want: defaultRetryAfter},
		{name: "falls back when unparseable", retryAfter: "later", want: defaultRetryAfter},
		{name: "falls back when not positive", retryAfter: "0", want: defaultRetryAfter},
		{name: "clamps a long wait", retryAfter: "600", want: maxRetryAfter},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := http.Header{}
			if tt.retryAfter != "" {
				header.Set("Retry-After", tt.retryAfter)
			}
			if got := retryAfterFrom(header); got != tt.want {
				t.Errorf("retryAfterFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleHTTPFailure(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		wantFatal   bool
		wantMessage string
	}{
		// A 429 used to fall through to the permanent branch and get logged as
		// a 5xx, which is what made the rate limit so hard to spot.
		{name: "rate limited is retryable", status: http.StatusTooManyRequests, wantFatal: false},
		{name: "server error is retryable", status: http.StatusInternalServerError, wantFatal: false},
		{name: "bad gateway is retryable", status: http.StatusBadGateway, wantFatal: false},
		{name: "unauthorized is permanent", status: http.StatusUnauthorized, wantFatal: true},
		{name: "bad request is permanent", status: http.StatusBadRequest, wantFatal: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handleHTTPFailure(response(tt.status, `{"ok":false,"error":"boom"}`, nil))
			if err.IsFatal() != tt.wantFatal {
				t.Errorf("handleHTTPFailure(%d).IsFatal() = %v, want %v", tt.status, err.IsFatal(), tt.wantFatal)
			}
			if !strings.Contains(err.Error(), "boom") {
				t.Errorf("handleHTTPFailure(%d) internal = %q, want it to include the response body", tt.status, err.Error())
			}
		})
	}
}

// The listing runs inside a synchronous OAuth callback. Retrying every page in
// turn is not enough on a workspace that stays rate limited across dozens of
// pages: without an overall budget the listing outlives the caller by minutes.
func TestClient_GetChannels_StopsWhenTheBudgetIsSpent(t *testing.T) {
	client, stub, _ := newTestClient(
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}],"response_metadata":{"next_cursor":"page2"}}`, nil),
		response(200, `{"ok":true,"channels":[{"id":"C2","name":"random"}],"response_metadata":{"next_cursor":"page3"}}`, nil),
		response(200, `{"ok":true,"channels":[{"id":"C3","name":"design"}],"response_metadata":{"next_cursor":"page4"}}`, nil),
	)

	// Spend the budget after the first page comes back.
	ctx, cancel := context.WithCancel(context.Background())
	stub.onCall = func(calls int) {
		if calls == 1 {
			cancel()
		}
	}
	defer cancel()

	got, err := client.GetChannels(ctx, &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v, want the pages fetched so far", err)
	}

	if want := []string{"general"}; !reflect.DeepEqual(channelNames(got), want) {
		t.Errorf("GetChannels() = %v, want %v", channelNames(got), want)
	}
	if stub.calls != 1 {
		t.Errorf("GetChannels() made %d requests, want it to stop after 1", stub.calls)
	}
}

// A caller that has gone away leaves nothing to return, and there is no point
// pretending the listing succeeded.
func TestClient_GetChannels_ErrorsWhenAbandonedWithNoResults(t *testing.T) {
	client, _, _ := newTestClient(
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}]}`, nil),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	got, err := client.GetChannels(ctx, &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err == nil {
		t.Fatal("GetChannels() error = nil, want an error for an abandoned listing")
	}
	if !isAbandoned(err) {
		t.Errorf("GetChannels() error = %v, want it to be an abandoned error", err)
	}
	if len(got) != 0 {
		t.Errorf("GetChannels() = %v, want no channels", got)
	}
}

// Backing off past the budget only delays a page nothing is left to fetch.
func TestClient_GetChannels_DoesNotBackOffPastTheBudget(t *testing.T) {
	client, stub, slept := newTestClient(
		response(200, `{"ok":true,"channels":[{"id":"C1","name":"general"}],"response_metadata":{"next_cursor":"page2"}}`, nil),
		rateLimitedResponse("1"),
		response(200, `{"ok":true,"channels":[{"id":"C2","name":"random"}]}`, nil),
	)

	// Spend the budget while the rate limited page is waiting to be retried.
	ctx, cancel := context.WithCancel(context.Background())
	stub.onCall = func(calls int) {
		if calls == 2 {
			cancel()
		}
	}
	defer cancel()

	got, err := client.GetChannels(ctx, &GetChannelsRequest{BearerToken: "xoxb-test"})
	if err != nil {
		t.Fatalf("GetChannels() unexpected error = %v", err)
	}

	if want := []string{"general"}; !reflect.DeepEqual(channelNames(got), want) {
		t.Errorf("GetChannels() = %v, want %v", channelNames(got), want)
	}
	// The wait is entered once and abandoned, never repeated.
	if len(*slept) != 1 {
		t.Errorf("backoff attempts = %d, want 1", len(*slept))
	}
	if stub.calls != 2 {
		t.Errorf("GetChannels() made %d requests, want it to stop after 2", stub.calls)
	}
}

// The budget must leave the caller room to receive the response.
func TestChannelListingBudget(t *testing.T) {
	if channelListingBudget < maxRetryAfter {
		t.Errorf("channelListingBudget = %v, want at least one retry to fit in it (%v)", channelListingBudget, maxRetryAfter)
	}
	if channelListingBudget > 30*time.Second {
		t.Errorf("channelListingBudget = %v, want it under the caller's own timeout", channelListingBudget)
	}
}

func TestSleep_ReturnsTheContextErrorWhenAbandoned(t *testing.T) {
	client := &Client{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := client.sleep(ctx, time.Hour); !errors.Is(err, context.Canceled) {
		t.Errorf("sleep() error = %v, want context.Canceled", err)
	}
}
