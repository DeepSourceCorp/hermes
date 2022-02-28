package provider

import (
	"context"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
)

type Provider interface {
	Send(context.Context, *domain.Notifier, string) (*domain.Message, domain.IError)
	Validate(context.Context, *domain.NotifierConfiguration) domain.IError
}

type IHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type BaseProvider struct {
	HTTPClient IHTTPClient
}
