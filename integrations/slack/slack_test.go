package slack

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// createMockServer returns a mock server with custom headers and response body for tests.
func createMockServer(ok bool, header int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		r := &response{}
		r.Ok = ok

		b, _ := json.Marshal(r)
		rw.WriteHeader(header)
		rw.Write(b)
	}))

	return server
}

// createMockClient returns a mock Client for tests.
func createMockClient() *Client {
	c := &Client{
		Client: http.DefaultClient,
		config: &Config{
			URI:   "https://api.slack.com",
			Token: "abc",
		},
	}

	return c
}

func TestNewClient(t *testing.T) {
	c := createMockClient()

	t.Run("config URI should be same as the value passed to NewClient()", func(t *testing.T) {
		got := NewClient("https://api.slack.com", "abc")
		if got.config.URI != "https://api.slack.com" {
			t.Errorf("NewClient() config.URI: Got = %v, want %v", got.config.URI, c.config.URI)
		}
	})
}

func TestSendMessage(t *testing.T) {
	c := createMockClient()

	server := createMockServer(true, 201)
	defer server.Close()
	c.config.URI = server.URL

	t.Run("status must be true while sending message", func(t *testing.T) {
		ctx := context.Background()
		status, _ := c.SendMessage(ctx, "general", "hello", nil)
		if status != true {
			t.Errorf("CreateIssue() status: Got = %v, want %v", status, "")
		}
	})
}

func TestSendMessageStatus(t *testing.T) {
	c := createMockClient()

	server := createMockServer(false, 201)
	defer server.Close()
	c.config.URI = server.URL

	t.Run("Status should be false when failed to send message", func(t *testing.T) {
		ctx := context.Background()
		status, _ := c.SendMessage(ctx, "general", "hello", nil)
		if status != false {
			t.Errorf("CreateIssue() err: Got = %v, want %v", status, "")
		}
	})
}
