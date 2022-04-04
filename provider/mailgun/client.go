package mailgun

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
)

type Client struct {
	HTTPClient provider.IHTTPClient
}

type SendMessageRequest struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Subject    string `json:"subject"`
	Text       string `json:"text"`
	Token      string `json:"-"`
	DomainName string `json:"domain"`
}

type SendMessageResponse struct {
	Ok bool `json:"ok"`
}

func (c *Client) SendMessage(request *SendMessageRequest) (interface{}, domain.IError) {
	data := url.Values{}
	data.Set("from", request.From)
	data.Set("to", request.To)
	data.Set("subject", request.Subject)
	data.Set("text", request.Text)

	uri := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", request.DomainName)

	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errFailedSendPermanent("failed to send request")
	}

	req.SetBasicAuth("api", request.Token)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errFailedSendTemporary("something went wrong while sending messsage to mailgun")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 399 {
		return nil, handleHTTPFailure(resp)
	}

	var response = new(SendMessageResponse)
	response.Ok = true
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
