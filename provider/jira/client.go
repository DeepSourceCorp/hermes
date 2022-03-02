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

type project struct {
	Key string `json:"key"`
}

type issueType struct {
	Name string `json:"name"`
}

type fields struct {
	Project     project                `json:"project"`
	IssueType   issueType              `json:"issuetype"`
	Summary     string                 `json:"summary"`
	Description map[string]interface{} `json:"description"`
}

type postIssueRequest struct {
	Fields      fields `json:"fields"`
	CloudID     string `json:"-"`
	BearerToken string `json:"-"`
}

type postIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func send(httpClient provider.IHTTPClient, request *postIssueRequest) (interface{}, domain.IError) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, errFailedSendPermanent("failed to encode request")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(postIssueURL, request.CloudID), &buf)
	if err != nil {
		return nil, errFailedSendTemporary("failed to send request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errFailedSendTemporary("something went wrong while creating the issue")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(postIssueResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errFailedSendPermanent("success but failed to parse response body")
	}
	return response, nil
}

func handleHTTPFailure(response *http.Response) domain.IError {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return errFailedSendPermanent("received non 2xx, failed to parse response")
	}

	if response.StatusCode > 500 {
		if err != nil {
			return errFailedSendTemporary(fmt.Sprintf("received 5xx, error=%s", string(b)))
		}
	}

	return errFailedSendPermanent(fmt.Sprintf("received 5xx, error=%s", string(b)))
}
