package mailgun

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
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

func TestMailgunSend(t *testing.T) {

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
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: false,
		},
		{
			name: "no secret in config",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: nil,
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "no sender",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "no receiver",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "no subject",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "no domain",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
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
				body: []byte(`{"text": "Hi Apollo!"}`),
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
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "token not set",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
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
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(``),
			},
			wantErr: true,
		},
		{
			name: "text not set",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text":""}`),
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
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"from": "mailgun@domain.com", "to": "apollo@deepsource.io", "subject": "Apollo + Mailgun", "domain": "secret",
						},
					},
				},
				body: []byte(`{"text": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &defaultMailgun{
				Client: &Client{HTTPClient: tt.fields.httpClient},
			}
			got, err := p.Send(tt.args.ctx, tt.args.notifier, tt.args.body)
			if tt.wantErr == false {
				if got.ID == "" {
					t.Errorf("defaultMailgun.Send() ID missing in payload")
				}
				if got.Ok == false {
					t.Errorf("defaultMailgun.Send() Ok == false")
				}
			}
			if err != nil != tt.wantErr {
				t.Errorf("defaultMailgun.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
