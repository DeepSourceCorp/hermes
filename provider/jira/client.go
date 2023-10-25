package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"

	log "github.com/sirupsen/logrus"
)

const (
	accessibleResourcesURL = "https://api.atlassian.com/oauth/token/accessible-resources"
	projectSearchURL       = "https://api.atlassian.com/ex/jira/%s/rest/api/3/project/search?expand=issueTypes"
	postIssueURL           = "https://api.atlassian.com/ex/jira/%s/rest/api/3/issue"
)

type Client struct {
	HTTPClient provider.IHTTPClient
}

type Project struct {
	Expand     string      `json:"expand"`
	Self       string      `json:"self"`
	ID         string      `json:"id"`
	Key        string      `json:"key"`
	IssueTypes []IssueType `json:"issueTypes"`
	Name       string      `json:"name"`
	AvatarUrls struct {
		Four8X48  string `json:"48x48"`
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
	} `json:"avatarUrls"`
	ProjectTypeKey string   `json:"projectTypeKey"`
	Simplified     bool     `json:"simplified"`
	Style          string   `json:"style"`
	IsPrivate      bool     `json:"isPrivate"`
	Properties     struct{} `json:"properties"`
}

type IssueType struct {
	Self           string `json:"self"`
	ID             string `json:"id"`
	Description    string `json:"description"`
	IconURL        string `json:"iconUrl"`
	Name           string `json:"name"`
	Subtask        bool   `json:"subtask"`
	AvatarID       int    `json:"avatarId,omitempty"`
	HierarchyLevel int    `json:"hierarchyLevel"`
}

type Reporter struct {
	ID string `json:"id,omitempty"`
}

type Fields struct {
	Project struct {
		Key string `json:"key"`
	} `json:"project"`
	IssueType struct {
		ID string `json:"id"`
	} `json:"issuetype"`
	Reporter    *Reporter              `json:"reporter,omitempty"`
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
		log.Errorf("jira: encoding request: %v", err)
		return nil, errFailedPermenant("failed to encode request")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(postIssueURL, request.CloudID), &buf)
	if err != nil {
		log.Errorf("jira: failed to send request: %v", err)
		return nil, errFailedTemporary("failed to send request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("jira: something went wrong while creating the issue: %v", err)
		return nil, errFailedTemporary("something went wrong while creating the issue")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	response := new(CreateIssueResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Errorf("jira: failed to parse response body: %v", err)
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

type AccessibleResourcesResponse []Site

func (c *Client) GetAccessibleResources(request *AccessibleResourcesRequest) (*AccessibleResourcesResponse, domain.IError) {
	req, err := http.NewRequest("GET", accessibleResourcesURL, http.NoBody)
	if err != nil {
		log.Errorf("jira: failed requesting accessible resources: %v", err)
		return nil, errFailedTemporary("failed to send request")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("jira: something went wrong requesting accessible resources: %v", err)
		return nil, errFailedTemporary("something went wrong")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	response := new(AccessibleResourcesResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Errorf("jira: failed decoding accessible resources response: %v", err)
		return nil, errFailedPermenant("success but failed to parse response body")
	}

	return response, nil
}

type GetProjectsRequest struct {
	BearerToken string
	CloudID     string
}

type GetProjectsResponse struct {
	Self       string    `json:"self"`
	MaxResults int       `json:"maxResults"`
	NextPage   string    `json:"nextPage"`
	StartAt    int       `json:"startAt"`
	Total      int       `json:"total"`
	IsLast     bool      `json:"isLast"`
	Values     []Project `json:"values"`
}

func (c *Client) GetProjects(request *GetProjectsRequest) ([]Project, domain.IError) {
	req, err := http.NewRequest("GET", fmt.Sprintf(projectSearchURL, request.CloudID), http.NoBody)
	if err != nil {
		log.Errorf("jira: failed requesting projects: %v", err)
		return nil, errFailedTemporary("failed to send request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	var projects []Project
	for {
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			log.Errorf("jira: something went wrong requesting projects: %v", err)
			return nil, errFailedTemporary("something went wrong")
		}

		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			return nil, handleHTTPFailure(resp)
		}

		response := new(GetProjectsResponse)
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			log.Errorf("jira: failed decoding projects response: %v", err)
			return nil, errFailedPermenant("success but failed to parse response body")
		}

		projects = append(projects, response.Values...)

		if response.IsLast {
			break
		}

		nextURL, err := url.Parse(response.NextPage)
		if err != nil {
			log.Errorf("jira: failed to parse next url in projects response: %v", err)
			break
		}
		req.URL = nextURL
	}
	return projects, nil
}

func handleHTTPFailure(response *http.Response) domain.IError {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("jira: received non 2xx, failed to parse response: %v", err)
		return errFailedPermenant("received non 2xx, failed to parse response")
	}

	// FIXME(SS): It returns the same error in the commented out code.
	// if response.StatusCode > http.StatusInternalServerError {
	// 	if err != nil {
	// 		return errFailedTemporary(fmt.Sprintf("received 5xx, error=%s", string(b)))
	// 	}
	// }

	log.Errorf("jira: received 5xx: %v", err)
	return errFailedPermenant(fmt.Sprintf("received 5xx, error=%s", string(b)))
}
