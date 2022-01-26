package jira

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Issue represents a Jira issue.
type Issue struct {
	Key    string       `json:"key"`
	Fields *IssueFields `json:"fields,omitempty"`
}

// IssueFields represents the fields for a Jira issue.
type IssueFields struct {
	Project     *Project   `json:"project,omitempty"`
	Summary     string     `json:"summary,omitempty"`
	Description string     `json:"description,omitempty"`
	IssueType   *IssueType `json:"issuetype,omitempty"`
	Status      string     `json:"status,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
}

// Project represents the project details for a Jira issue.
type Project struct {
	Key string `json:"key,omitempty"`
}

// IssueType represents the issue type for a Jira issue.
type IssueType struct {
	Name string `json:"name,omitempty"`
}

// Config stores client configuration.
type Config struct {
	URI        string
	Username   string
	Token      string
	APIVersion string
}

// Client represents a Jira client.
type Client struct {
	Client *http.Client
	config *Config
}

// NewClient returns a Jira client.
func NewClient(uri, username, token, apiVersion string) *Client {
	return &Client{
		Client: http.DefaultClient,
		config: &Config{
			URI:        uri,
			Username:   username,
			Token:      token,
			APIVersion: apiVersion,
		},
	}
}

// CreateIssue creates a Jira issue and returns the key for the issue, along with an error.
func (client *Client) CreateIssue(issue Issue) (key string, err error) {
	b, err := json.Marshal(issue)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", client.config.URI+"/rest/api/"+client.config.APIVersion+"/issue/", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	// set request headers
	req.SetBasicAuth(client.config.Username, client.config.Token)
	req.Header.Set("Content-Type", "application/json")

	// perform request
	resp, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}

	// read response body
	respBody, _ := io.ReadAll(resp.Body)
	var parsed map[string]interface{}
	json.Unmarshal([]byte(string(respBody)), &parsed)

	// set issue key
	issue.Key = parsed["key"].(string)

	if resp.StatusCode == 201 {
		return issue.Key, nil
	} else {
		return "", errors.New("Error:" + resp.Status)
	}
}
