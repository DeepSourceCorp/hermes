package smtp

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type defaultSMTP struct {
	Client *Client
}

const ProviderType = domain.ProviderType("smtp")

func NewSMTPProvider(httpClient provider.IHTTPClient) provider.Provider {
	return &defaultSMTP{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *defaultSMTP) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
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
		FromEmail:    opts.FromEmail,
		FromPassword: opts.FromPassword,
		ToEmail:      opts.ToEmail,
		SMTPHost:     opts.SMTPHost,
		SMTPPort:     opts.SMTPPort,
		Subject:      opts.Subject,
		Message:      payload.Message,
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
	FromEmail    string   `mapstructure:"from_email"`
	FromPassword string   `mapstructure:"from_password"`
	ToEmail      []string `mapstructure:"to_email"`
	SMTPHost     string   `mapstructure:"smtp_host"`
	SMTPPort     string   `mapstructure:"smtp_port"`
	Subject      string   `mapstructure:"subject"`
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
	if o.FromEmail == "" || o.FromPassword == "" {
		return errFailedOptsValidation("sender email and password is empty")
	}
	if o.ToEmail == nil {
		return errFailedOptsValidation("reciever mail is empty")
	}
	if o.SMTPHost == "" || o.SMTPPort == "" {
		return errFailedOptsValidation("SMTP host and port is empty")
	}
	if o.Subject == "" {
		return errFailedOptsValidation("subject is empty")
	}
	return nil
}

type Payload struct {
	Message string `json:"message"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		p.Message = string(body)
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Message == "" {
		return errFailedBodyValidation("message is empty")
	}
	return nil
}
