package linear

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type defaultLinear struct {
	Client *Client
}

const ProviderType = domain.ProviderType("linear")

func NewLinearProvider(httpClient provider.IHTTPClient) provider.Provider {
	return &defaultLinear{
		Client: &Client{HTTPClient: httpClient},
	}
}

func (p *defaultLinear) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
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
		TeamID:      opts.TeamID,
		Description: payload.Description,
		Title:       opts.Title,
		BearerToken: opts.Secret.Token,
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
	TeamID string `mapstructure:"teamId"`
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
	if o.TeamID == "" {
		return errFailedOptsValidation("teamID is empty")
	}
	return nil
}

type Payload struct {
	Description string `json:"description"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		body := string(body)
		p.Description = strings.TrimSuffix(body, "\n")
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Description == "" {
		return errFailedBodyValidation("description is empty")
	}
	return nil
}
