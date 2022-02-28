package slack

import (
	"context"
	"errors"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
)

type providerSlack struct {
	provider.BaseProvider
}

func NewSlackProvider(httpClient *http.Client) provider.Provider {
	return &providerSlack{
		provider.BaseProvider{
			HTTPClient: httpClient,
		},
	}
}

type Opts struct {
	Token   string
	Channel string
}

func (p *providerSlack) Send(ctx context.Context, notifier *domain.Notifier, body string) (*domain.Message, domain.IError) {
	opts, err := newOpts(notifier.Config)
	if err != nil {
		return nil, errSlackOptsParseFail(err.Error())
	}

	request := &postMessageRequest{
		Channel: opts.Channel,
		Token:   opts.Token,
		Text:    body,
	}

	c := client{
		httpClient: p.BaseProvider.HTTPClient,
	}

	response, err := c.SendMessage(request)
	if err != nil {
		return nil, errSlackErr(err.Error())
	}

	return &domain.Message{
		ID:               ksuid.New().String(),
		Ok:               true,
		ProviderResponse: response,
		Body:             body,
	}, nil
}

func newOpts(config *domain.NotifierConfiguration) (*Opts, error) {
	var opts = new(Opts)
	if err := mapstructure.Decode(config.Opts, opts); err != nil {
		return nil, err
	}
	opts.Token = config.Secret.Token
	if opts.Channel == "" || opts.Token == "" {
		return nil, errors.New("error parsing slack configuration")
	}
	return opts, nil
}

func (*providerSlack) Validate(ctx context.Context, config *domain.NotifierConfiguration) domain.IError {
	return nil
}
