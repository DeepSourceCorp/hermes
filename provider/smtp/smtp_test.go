package smtp

import (
	"bytes"
	"context"
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

func TestSMTPSend(t *testing.T) {

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
							"from_email": "test@domain.com", "from_password": "abcd", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_host": "localhost", "smtp_port": "1025",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
			},
			wantErr: false,
		},
		{
			name: "no password",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"from_email": "test@domain.com", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_host": "localhost", "smtp_port": "1025",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
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
						Opts: map[string]interface{}{
							"to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_host": "localhost", "smtp_port": "1025",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
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
						Opts: map[string]interface{}{
							"from_email": "test@domain.com", "from_password": "abcd", "subject": "Apollo + SMTP", "smtp_host": "localhost", "smtp_port": "1025",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
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
						Opts: map[string]interface{}{
							"from_email": "test@domain.com", "from_password": "abcd", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "smtp_host": "localhost", "smtp_port": "1025",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "no smtp host",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"from_email": "test@domain.com", "from_password": "abcd", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_port": "1025",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
			},
			wantErr: true,
		},
		{
			name: "no smtp port",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"from_email": "test@domain.com", "from_password": "abcd", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_host": "localhost",
						},
					},
				},
				body: []byte(`{"message": "Hi Apollo!"}`),
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
				body: []byte(`{"message": "Hi Apollo!"}`),
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
				body: []byte(`{"message": "Hi Apollo!"}`),
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
							"from_email": "test@domain.com", "from_password": "abcd", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_host": "localhost",
						},
					},
				},
				body: []byte(``),
			},
			wantErr: true,
		},
		{
			name: "message not set",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Opts: map[string]interface{}{
							"from_email": "test@domain.com", "from_password": "abcd", "to_email": []string{"apollo@deepsource.io", "test@deepsource.io"}, "subject": "Apollo + SMTP", "smtp_host": "localhost",
						},
					},
				},
				body: []byte(`{"message":""}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &defaultSMTP{
				Client: &Client{HTTPClient: tt.fields.httpClient},
			}
			got, err := p.Send(tt.args.ctx, tt.args.notifier, tt.args.body)
			if tt.wantErr == false {
				if got.Ok == false {
					t.Errorf("defaultSMTP.Send() Ok == false")
				}
			}
			if err != nil != tt.wantErr {
				t.Errorf("defaultSMTP.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
