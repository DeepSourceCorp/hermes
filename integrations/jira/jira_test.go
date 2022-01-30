package jira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// createMockServer returns a mock server with custom headers and response body for tests.
func createMockServer(id string, header int) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		r := map[string]string{
			"key": string(id),
		}

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
			URI:        "https://jira.atlassian.net",
			Token:      "abc",
			Username:   "test",
			APIVersion: "latest",
		},
	}

	return c
}

// createMockIssue returns a mock Issue for tests.
func createMockIssue() Issue {
	i := Issue{
		Key: "1",
		Fields: &IssueFields{
			Project: &Project{
				Key: "Test",
			},
			Summary:     "Test Issue",
			Description: "Test Description",
		},
	}

	return i
}

func TestNewClient(t *testing.T) {
	c := createMockClient()

	t.Run("config URI should be same as the value passed to NewClient()", func(t *testing.T) {
		got := NewClient("https://jira.atlassian.net", "test", "abc", "latest")
		if got.config.URI != "https://jira.atlassian.net" {
			t.Errorf("NewClient() config.URI: Got = %v, want %v", got.config.URI, c.config.URI)
		}
	})
}

func TestCreateIssue(t *testing.T) {
	c := createMockClient()
	i := createMockIssue()

	server := createMockServer("1", 201)
	defer server.Close()
	c.config.URI = server.URL

	t.Run("ID should be 1", func(t *testing.T) {
		id, _ := c.CreateIssue(i)
		if id != "1" {
			t.Errorf("CreateIssue() id: Got = %v, want %v", id, "1")
		}
	})
}

func TestCreateIssueWithoutServer(t *testing.T) {
	c := createMockClient()
	i := createMockIssue()

	t.Run("ID should be blank while creating an issue without a server", func(t *testing.T) {
		id, _ := c.CreateIssue(i)
		if id != "" {
			t.Errorf("CreateIssue() id: Got = %v, want %v", id, "")
		}
	})
}

func TestCreateIssue404(t *testing.T) {
	c := createMockClient()
	i := createMockIssue()

	server := createMockServer("1", 404)
	defer server.Close()
	c.config.URI = server.URL

	t.Run("error should be 404 Not Found", func(t *testing.T) {
		id, err := c.CreateIssue(i)
		if err.Error() != "Error: 404 Not Found" {
			t.Errorf("CreateIssue() err: Got = %v, want %v", id, "")
		}
	})
}
