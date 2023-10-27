package slack

import (
	"context"
	"encoding/json"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"

	log "github.com/sirupsen/logrus"
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

func (p *defaultSlack) Send(_ context.Context, notifier *domain.Notifier, body []byte) (*domain.Message, domain.IError) {
	// Extract and validate the payload.
	var payload = new(Payload)
	if err := payload.Extract(body); err != nil {
		log.Errorf("slack: failed extracting payload: %v", err)
		return nil, err
	}
	if err := payload.Validate(); err != nil {
		log.Errorf("slack: failed validating payload: %v", err)
		return nil, err
	}

	// Extract and validate the configuration.
	var opts = new(Opts)
	if err := opts.Extract(notifier.Config); err != nil {
		log.Errorf("slack: failed extracting options: %v", err)
		return nil, err
	}

	if err := opts.Validate(); err != nil {
		log.Errorf("slack: failed validating options: %v", err)
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
		log.Errorf("slack: failed sending message: %v", err)
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
		log.Errorf("slack: failed decoding options: %v", err)
		return errFailedOptsValidation("failed to decode configuration")
	}
	o.Secret = c.Secret
	return nil
}

func (o *Opts) Validate() domain.IError {
	if o.Secret == nil || o.Secret.Token == "" {
		log.Errorf("slack: secret not defined in configuration")
		return errFailedOptsValidation("secret not defined in configuration")
	}
	if o.Channel == "" {
		log.Errorf("slack: channel not provided")
		return errFailedOptsValidation("channel is empty")
	}
	return nil
}

type Payload struct {
	Text   string      `json:"text"`
	Blocks interface{} `json:"blocks"`
}

func (p *Payload) Extract(body []byte) domain.IError {
	if err := json.Unmarshal(body, p); err != nil {
		log.Errorf("slack: failed to unmarshal: %v", err)
		p.Text = string(body)
	}
	return nil
}

func (p *Payload) Validate() domain.IError {
	if p.Blocks == nil && p.Text == "" {
		log.Errorf("slack: failed to validate")
		return errFailedBodyValidation("blocks and text is empty")
	}
	return nil
}

func (p *defaultSlack) GetOptValues(_ context.Context, secret *domain.NotifierSecret) (map[string]interface{}, error) {
	request := &GetChannelsRequest{
		BearerToken: secret.Token,
	}
	channels, err := p.Client.GetChannels(request)
	if err != nil {
		log.Errorf("slack: failed to get channels : %v", err)
		return nil, err
	}

	return map[string]interface{}{"channel": channels}, nil
}
