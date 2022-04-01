package github

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type defaultGithub struct {
	Client *Client
}

const ProviderType = domain.ProviderType("github")

func NewGithubProvider(httpClient provider.IHTTPClient) provider.Provider {
	return &defaultGithub{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *defaultGithub) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
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
		Repo:  opts.Repo,
		Owner: opts.Owner,
		Title: opts.Title,
		Token: opts.Secret.Token,
		Body:  payload.Body,
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
	Secret *domain.NotifierSecret
	Repo   string `mapstructure:"repo"`
	Owner  string `mapstructure:"owner"`
	Title  string `mapstructure:"title"`
}

func (o *Opts) Extract(c *domain.NotifierConfiguration) domain.IError {
	if c == nil {
		return errFailedOptsValidation("notifier config empty")
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
	if o.Repo == "" || o.Owner == "" {
		return errFailedOptsValidation("repo and owner is emtpy")
	}
	if o.Title == "" {
		return errFailedOptsValidation("title is emtpy")
	}
	return nil
}

type Payload struct {
	Body string `json:"body"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		p.Body = string(body)
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Body == "" {
		return errFailedBodyValidation("body is empty")
	}
	return nil
}
