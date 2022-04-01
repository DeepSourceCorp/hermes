package mailgun

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type defaultMailgun struct {
	Client *Client
}

const ProviderType = domain.ProviderType("mailgun")

func NewMailgunProvider(httpClient provider.IHTTPClient) provider.Provider {
	return &defaultMailgun{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *defaultMailgun) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
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
		From:       opts.From,
		To:         opts.To,
		Subject:    opts.Subject,
		Text:       payload.Text,
		Token:      opts.Secret.Token,
		DomainName: opts.DomainName,
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
	Secret     *domain.NotifierSecret
	From       string `mapstructure:"from"`
	To         string `mapstructure:"to"`
	Subject    string `mapstructure:"subject"`
	DomainName string `mapstructure:"domain"`
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
	if o.From == "" || o.To == "" {
		return errFailedOptsValidation("sender and receipient details is empty")
	}
	if o.Subject == "" {
		return errFailedOptsValidation("subject is empty")
	}
	if o.DomainName == "" {
		return errFailedOptsValidation("domain is empty")
	}

	return nil
}

type Payload struct {
	Text string `json:"text"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		p.Text = string(body)
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Text == "" {
		return errFailedBodyValidation("text is empty")
	}
	return nil
}
