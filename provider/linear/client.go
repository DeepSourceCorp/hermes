package linear

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
)

const postMessageURL = "https://api.linear.app/graphql"

type Client struct {
	HTTPClient provider.IHTTPClient
}

type SendMessageRequest struct {
	TeamID      string `json:"teamId"`
	Description string `json:"description,omitempty"`
	Title       string `json:"title,omitempty"`
	BearerToken string `json:"-"`
}

type SendMessageResponse struct {
	Ok bool `json:"ok"`
}

func (c *Client) SendMessage(request *SendMessageRequest) (interface{}, domain.IError) {
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
		"title":       request.Title,
		"description": request.Description,
		"teamID":      request.TeamID,
	}

	// parse query template
	t := template.Must(template.New("json").Parse(queryTemplate))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, queryData); err != nil {
		return nil, errFailedSendPermanent("failed to encode request")
	}

	query := map[string]string{
		"query": buf.String(),
	}

	b, err := json.Marshal(query)
	if err != nil {
		return nil, errFailedSendPermanent("failed to encode request")
	}

	req, err := http.NewRequest("POST", postMessageURL, bytes.NewBuffer(b))
	if err != nil {
		return nil, errFailedSendPermanent("failed to send request")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedSendTemporary("something went wrong while sending messsage to linear")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(SendMessageResponse)
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
