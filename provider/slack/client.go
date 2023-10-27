package slack

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
		log.Errorf("slack: failed encoding request: %v", err)
		return nil, errFailedSendPermanent(err.Error())
	}

	req, err := http.NewRequest("POST", postMessageURL, &buf)
	if err != nil {
		log.Errorf("slack: sending request: %v", err)
		return nil, errFailedSendPermanent(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("slack: something went wrong sending request: %v", err)
		return nil, errFailedSendTemporary("something went wrong while sending messsage to slack")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(SendMessageResponse)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Errorf("slack: failed decoding response: %v", err)
		return nil, errFailedSendPermanent(err.Error())
	}
	return response, nil
}

func handleHTTPFailure(response *http.Response) domain.IError {
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("slack: failed reading response body: %v", err)
		return errFailedSendPermanent(err.Error())
	}

	if response.StatusCode > 500 {
		if err != nil {
			log.Errorf("slack: failed with 5xx response code: %v", err)
			return errFailedSendTemporary(fmt.Sprintf("received 5xx, error=%s", string(b)))
		}
	}

	log.Errorf("slack: failed with 5xx response code: %v", err)
	return errFailedSendPermanent(fmt.Sprintf("received 5xx, error=%s", string(b)))
}

const getChannelsURL = "https://slack.com/api/conversations.list?types=public_channel,private_channel&exclude_archived=true&limit=1000"

type GetChannelsRequest struct {
	BearerToken string `json:"_"`
}

type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type GetChannelsResponse struct {
	Ok               bool             `json:"ok"`
	Channels         []Channel        `json:"channels"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
}

func (c *Client) getChannelsPage(request *GetChannelsRequest, cursor string) (*GetChannelsResponse, domain.IError) {
	var response = new(GetChannelsResponse)

	requestUrl := getChannelsURL
	if cursor != "" {
		requestUrl += "&cursor=" + url.QueryEscape(cursor)
	}

	req, err := http.NewRequest("GET", requestUrl, http.NoBody)
	if err != nil {
		log.Errorf("slack: failed creating request for options: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", request.BearerToken))
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Errorf("slack: failed sending request for options: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		log.Errorf("slack: Non 2xx response while fetching options: %v", err)
		return response, handleHTTPFailure(resp)
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		log.Errorf("slack: Non 2xx response while fetching options: %v", err)
		return response, errFailedOptsFetch(err.Error())
	}

	return response, nil
}

func (c *Client) GetChannels(request *GetChannelsRequest) ([]map[string]string, domain.IError) {
	channels := []map[string]string{}
	cursor := ""

	for {
		response, err := c.getChannelsPage(request, cursor)
		if err != nil {
			log.Errorf("slack: Error fetching page %v: %v", cursor, err)
			return channels, err
		}

		for _, v := range response.Channels {
			channels = append(channels, map[string]string{
				"id":   v.ID,
				"name": v.Name,
			})
		}

		cursor = response.ResponseMetadata.NextCursor
		if cursor == "" {
			break
		}
	}

	return channels, nil
}
