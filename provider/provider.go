package provider

import (
	"context"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
)

type Provider interface {
	Send(context.Context, *domain.Notifier, []byte) (*domain.Message, domain.IError)
	GetOptValues(context.Context, *domain.NotifierSecret) (*map[string]interface{}, error)
}

type IHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
