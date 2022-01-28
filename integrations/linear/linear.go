package linear

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"text/template"
)

// Issue represents a Linear issue.
type Issue struct {
	ID          string `json:"id"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

// Config stores client configuration.
type Config struct {
	URI    string
	Token  string
	TeamID string
}

// Client represents a Linear client.
type Client struct {
	Client *http.Client
	config *Config
}

// NewClient returns a Linear client.
func NewClient(uri, token, teamID string) *Client {
	return &Client{
		Client: http.DefaultClient,
		config: &Config{
			URI:    uri,
			Token:  token,
			TeamID: teamID,
		},
	}
}

// CreateIssue creates a Linear issue and returns the ID for the issue, along with an error.
func (client *Client) CreateIssue(issue Issue) (id string, err error) {
	queryTemplate := `mutation IssueCreate {
		issueCreate(
			input: {
				title: "{{.title}}"
				description: "{{.description}}"
				teamId: "{{.teamID}}"
			}
		) {
			success
			issue {
				id
				title
			}
		}
	}`

	queryData := map[string]interface{}{
		"title":       issue.Title,
		"description": issue.Description,
		"teamID":      client.config.TeamID,
	}

	// parse query template
	t := template.Must(template.New("json").Parse(queryTemplate))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, queryData); err != nil {
		return "", err
	}

	query := map[string]string{
		"query": buf.String(),
	}

	b, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", client.config.URI, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	// set request headers
	req.Header.Set("Authorization", client.config.Token)
	req.Header.Set("Content-Type", "application/json")

	// perform request
	resp, err := client.Client.Do(req)
	if err != nil {
		return "", err
	}

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed map[string]interface{}
	json.Unmarshal([]byte(string(respBody)), &parsed)

	// set issue ID
	issue.ID = parsed["data"].(map[string]interface{})["issueCreate"].(map[string]interface{})["issue"].(map[string]interface{})["id"].(string)

	if resp.StatusCode == 200 {
		return issue.ID, nil
	}

	return "", errors.New("Error: " + resp.Status)
}
