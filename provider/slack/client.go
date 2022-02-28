package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deepsourcelabs/hermes/provider"
)

type client struct {
	httpClient provider.IHTTPClient
}

const postMessageURL = "https://slack.com/api/chat.postMessage"

type postMessageRequest struct {
	Channel string `json:"channel"`
	Blocks  string `json:"blocks,omitempty"`
	Text    string `json:"text,omitempty"`
	Token   string `json:"-"`
}

type postMessageResponse struct {
	Ok bool `json:"ok"`
}

func (c *client) SendMessage(request *postMessageRequest) (interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", postMessageURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.Token))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to send message to Slack, error=%v", err)
	}

	response := new(postMessageResponse)

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, fmt.Errorf("failed to send message to Slack, error=%v", err)
	}

	return response, nil
}
