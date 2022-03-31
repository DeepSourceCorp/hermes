package discord

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/stretchr/testify/assert"
)

type mockHttp struct{}

func (*mockHttp) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		Body:       io.NopCloser(bytes.NewReader([]byte("{\"ok\":true}"))),
		StatusCode: http.StatusOK,
	}, nil
}

type errHTTP struct{}

func (*errHTTP) Do(_ *http.Request) (*http.Response, error) {
	return nil, errors.New("test")
}

func TestDiscordSend(t *testing.T) {

	type fields struct {
		httpClient provider.IHTTPClient
	}
	type args struct {
		ctx      context.Context
		notifier *domain.Notifier
		body     []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Message
		wantErr bool
	}{
		{
			name: "valid request",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"webhook": "https://discord.com/api/webhooks/958",
						},
					},
				},
				body: []byte(`{"content": "Hi Apollo!"}`),
			},
			wantErr: false,
		},
		{
			name: "no webhook URI",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"webhook": "",
						},
					},
				},
				body: []byte(`{"description": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "config not set",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: nil,
				},
				body: []byte(`{"description": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "opts not set",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{},
				},
				body: []byte(`{"description": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "empty body",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"webhook": "https://discord.com/api/webhooks/958",
						},
					},
				},
				body: []byte(``),
			},
			wantErr: true,
		},
		{
			name: "http errors",
			fields: fields{
				httpClient: new(errHTTP),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"webhook": "https://discord.com/api/webhooks/958",
						},
					},
				},
				body: []byte(`{"description": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &defaultDiscord{
				Client: &Client{HTTPClient: tt.fields.httpClient},
			}
			got, err := p.Send(tt.args.ctx, tt.args.notifier, tt.args.body)
			if tt.wantErr == false {
				if got.ID == "" {
					t.Errorf("defaultDiscord.Send() ID missing in payload")
				}
				if got.Ok == false {
					t.Errorf("defaultDiscord.Send() Ok == false")
				}
			}
			if err != nil != tt.wantErr {
				t.Errorf("defaultDiscord.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpts_Extract(t *testing.T) {
	nc := &domain.NotifierConfiguration{
		Opts: map[string]interface{}{
			"webhook": "",
		},
	}

	got := new(Opts)
	got.Extract(nc)
	// validate should return an error when webhook URI is empty
	err := got.Validate()
	assert.NotNil(t, err)

	// extract should return an error when notifier config is empty
	err = got.Extract(nil)
	assert.NotNil(t, err)
}
