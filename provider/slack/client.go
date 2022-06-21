package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"

	log "github.com/sirupsen/logrus"
)

const postMessageURL = "https://slack.com/api/chat.postMessage"

type Client struct {
	HTTPClient provider.IHTTPClient
}

type SendMessageRequest struct {
	Channel     string      `json:"channel"`
	Blocks      interface{} `json:"blocks,omitempty"`
	Text        string      `json:"text,omitempty"`
	BearerToken string      `json:"-"`
}

type SendMessageResponse struct {
	Ok bool `json:"ok"`
}

func (c *Client) SendMessage(request *SendMessageRequest) (*SendMessageResponse, domain.IError) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		log.Error(err.Error())
		return nil, errFailedSendPermanent(err.Error())
	}

	req, err := http.NewRequest("POST", postMessageURL, &buf)
	if err != nil {
		return nil, errFailedSendPermanent(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedSendTemporary("something went wrong while sending messsage to slack")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	response := new(SendMessageResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errFailedSendPermanent(err.Error())
	}
	return response, nil
}

func handleHTTPFailure(response *http.Response) domain.IError {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return errFailedSendPermanent(err.Error())
	}

	if response.StatusCode > 500 {
		if err != nil {
			return errFailedSendTemporary(fmt.Sprintf("received 5xx, error=%s", string(b)))
		}
	}

	return errFailedSendPermanent(fmt.Sprintf("received 5xx, error=%s", string(b)))
}

const getChannelsURL = "https://slack.com/api/conversations.list?exclude_archived=true&limit=1000"

type GetChannelsRequest struct {
	BearerToken string `json:"_"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetChannelsResponse struct {
	Ok       bool      `json:"ok"`
	Channels []Channel `json:"channels"`
}

func (c *Client) GetChannels(request *GetChannelsRequest) (*GetChannelsResponse, domain.IError) {
	req, err := http.NewRequest("GET", getChannelsURL, nil)
	if err != nil {
		return nil, errFailedOptsFetch(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedOptsFetch(err.Error())
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	response := new(GetChannelsResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, errFailedOptsFetch(err.Error())
	}

	return response, nil
}
