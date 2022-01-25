package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// request represents a typical HTTP request sent to the Slack API.
type request struct {
	Channel     string      `json:"channel"`
	Text        string      `json:"text"`
	Attachments interface{} `json:"attachments"`
}

// Attachment represents a basic Slack Attachment
type Attachment struct {
	Fallback   string             `json:"fallback"`
	Color      string             `json:"color"`
	Pretext    string             `json:"pretext"`
	AuthorName string             `json:"author_name"`
	AuthorLink string             `json:"author_link"`
	AuthorIcon string             `json:"author_icon"`
	Title      string             `json:"title"`
	TitleLink  string             `json:"title_link"`
	Text       string             `json:"text"`
	Fields     []*AttachmentField `json:"fields"`
	ImageUrl   string             `json:"image_url"`
	ThumbUrl   string             `json:"thumb_url"`
	Footer     string             `json:"footer"`
	FooterIcon string             `json:"footer_icon"`
	Ts         int64              `json:"ts"`
}

// AttachmentField represents Slack fields.
type AttachmentField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// response represents the returned HTTP response.
type response struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

// Config stores client configuration.
type Config struct {
	URI   string
	Token string
}

// Client contains a http.Client along with config.
type Client struct {
	Client *http.Client
	config *Config
}

// NewClient returns a Client.
func NewClient(uri string, token string) *Client {
	return &Client{
		Client: http.DefaultClient,
		config: &Config{
			URI:   uri,
			Token: token,
		},
	}
}

// Sends message to the specified channel.
func (client *Client) SendMessage(ctx context.Context, channel, text string, attachments []*Attachment) (bool, error) {
	req := &request{
		Channel:     channel,
		Text:        text,
		Attachments: attachments,
	}

	// marshal request into JSON
	b, err := json.Marshal(req)
	if err != nil {
		return false, err
	}

	status, err := client.Request(ctx, b)
	if err != nil {
		return status, err
	}

	return status, nil
}

// Sends a HTTP request to the URI and decodes response.
func (client *Client) Request(ctx context.Context, b []byte) (bool, error) {
	// create new request
	req, err := http.NewRequest(http.MethodPost, client.config.URI, bytes.NewBuffer(b))
	if err != nil {
		return false, err
	}

	// add headers
	_ = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+client.config.Token)
	req.Header.Set("Content-type", "application/json; charset=utf-8")

	// send HTTP request
	resp, err := client.Client.Do(req)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	// decode response
	var res response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}

	// check response and return error
	if !res.Ok {
		return false, err
	}

	return res.Ok, nil
}
