package discord

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type defaultDiscord struct {
	Client *Client
}

const ProviderType = domain.ProviderType("discord")

func NewDiscordProvider(httpClient provider.IHTTPClient) provider.Provider {
	return &defaultDiscord{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *defaultDiscord) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
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
		Content:    payload.Content,
		WebhookURI: opts.WebhookURI,
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

type Payload struct {
	Content string `json:"content"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		p.Content = string(body)
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Content == "" {
		return errFailedBodyValidation("content is empty")
	}
	return nil
}

type Opts struct {
	WebhookURI string `mapstructure:"webhook"`
}

func (o *Opts) Extract(c *domain.NotifierConfiguration) domain.IError {
	if c == nil {
		return errFailedOptsValidation("notifier config empty")
	}
	if err := mapstructure.Decode(c.Opts, o); err != nil {
		return errFailedOptsValidation("failed to decode configuration")
	}
	return nil
}

func (o *Opts) Validate() domain.IError {
	if o.WebhookURI == "" {
		return errFailedOptsValidation("webhook URI is empty")
	}
	return nil
}
