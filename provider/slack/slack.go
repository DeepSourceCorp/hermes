package slack

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type defaultSlack struct {
	Client *Client
}

const ProviderType = domain.ProviderType("slack")

func NewSlackProvider(httpClient provider.IHTTPClient) provider.Provider {
	return &defaultSlack{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *defaultSlack) Send(ctx context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
	// Extract and validate the payload.
	var payload = new(Payload)
	if err := payload.Extract(body); err != nil {
		return nil, err
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Extract and validate the configuration.
	var opts = new(Opts)
	if err := opts.Extract(notifier.Config); err != nil {
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		return nil, err
	}

	request := &SendMessageRequest{
		Channel:     opts.Channel,
		BearerToken: opts.Secret.Token,
		Text:        payload.Text,
		Blocks:      payload.Blocks,
	}

	response, err := p.Client.SendMessage(request)
	if err != nil {
		return nil, err
	}

	return &domain.Message{
		ID:               ksuid.New().String(),
		Ok:               true,
		Payload:          payload,
		ProviderResponse: response,
	}, nil
}

type Opts struct {
	Secret  *domain.NotifierSecret
	Channel string `mapstructure:"channel"`
}

func (o *Opts) Extract(c *domain.NotifierConfiguration) domain.IError {
	if c == nil {
		return errFailedOptsValidation("notifier config emtpy")
	}
	if err := mapstructure.Decode(c.Opts, o); err != nil {
		return errFailedOptsValidation("failed to decode configuration")
	}
	o.Secret = c.Secret
	return nil
}

func (o *Opts) Validate() domain.IError {
	if o.Secret == nil || o.Secret.Token == "" {
		return errFailedOptsValidation("secret not defined in configuration")
	}
	if o.Channel == "" {
		return errFailedOptsValidation("channel is emtpy")
	}
	return nil
}

type Payload struct {
	Text   string      `json:"text"`
	Blocks interface{} `json:"blocks"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		p.Text = string(body)
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Blocks == nil && p.Text == "" {
		return errFailedBodyValidation("blocks and text is empty")
	}
	return nil
}
