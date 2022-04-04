package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
)

type Client struct {
	HTTPClient provider.IHTTPClient
}

type SendMessageRequest struct {
	Content    string `json:"content"`
	Username   string `json:"username"`
	WebhookURI string `json:"webhook"`
}

type SendMessageResponse struct {
	Ok bool `json:"ok"`
}

func (c *Client) SendMessage(request *SendMessageRequest) (interface{}, domain.IError) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return nil, errFailedSendPermanent("failed to encode request")
	}

	req, err := http.NewRequest("POST", request.WebhookURI, &buf)
	if err != nil {
		return nil, errFailedSendPermanent("failed to send request")
	}

	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedSendTemporary("something went wrong while sending messsage to discord")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(SendMessageResponse)
	if resp.StatusCode == 204 {
		response.Ok = true
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
