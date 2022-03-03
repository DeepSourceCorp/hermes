package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
)

const postIssueURL = "https://api.atlassian.com/ex/jira/%s/rest/api/3/issue"
const accessibleResourcesURL = "https://api.atlassian.com/oauth/token/accessible-resources"

type Client struct {
	HTTPClient provider.IHTTPClient
}

type Project struct {
	Key string `json:"key"`
}

type IssueType struct {
	Name string `json:"name"`
}

type Fields struct {
	Project     Project                `json:"project"`
	IssueType   IssueType              `json:"issuetype"`
	Summary     string                 `json:"summary"`
	Description map[string]interface{} `json:"description"`
}

type CreateIssueRequest struct {
	Fields      Fields `json:"fields"`
	CloudID     string `json:"-"`
	BearerToken string `json:"-"`
}

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func (c *Client) CreateIssue(request *CreateIssueRequest) (*CreateIssueResponse, domain.IError) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, errFailedPermenant("failed to encode request")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(postIssueURL, request.CloudID), &buf)
	if err != nil {
		return nil, errFailedTemporary("failed to send request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedTemporary("something went wrong while creating the issue")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(CreateIssueResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errFailedPermenant("success but failed to parse response body")
	}
	return response, nil
}

type Site struct {
	ID        string   `json:"id"`
	URL       string   `json:"url"`
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	AvatarURL string   `json:"avatarUrl"`
}

type AccessibleResourcesRequest struct {
	BearerToken string
}

type AccessibleResourcesResponse struct {
	Sites []Site
}

func (c *Client) GetAccessibleResources(request *AccessibleResourcesRequest) (*AccessibleResourcesResponse, domain.IError) {
	req, err := http.NewRequest("GET", accessibleResourcesURL, nil)
	if err != nil {
		return nil, errFailedTemporary("failed to send request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedTemporary("something went wrong")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(AccessibleResourcesResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errFailedPermenant("success but failed to parse response body")
	}
	return response, nil
}

func handleHTTPFailure(response *http.Response) domain.IError {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return errFailedPermenant("received non 2xx, failed to parse response")
	}

	if response.StatusCode > 500 {
		if err != nil {
			return errFailedTemporary(fmt.Sprintf("received 5xx, error=%s", string(b)))
		}
	}

	return errFailedPermenant(fmt.Sprintf("received 5xx, error=%s", string(b)))
}
