package jira

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
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

func Test_jiraSimple_Send(t *testing.T) {
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
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
			},
			wantErr: false,
		},
		{
			name: "no project_key",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
			},
			wantErr: true,
		},
		{
			name: "no issue_type",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"project_key": "abc",
						},
					},
				},
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
			},
			wantErr: true,
		},
		{
			name: "no summary in body",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"description":{"x":"y"}}`),
			},
			wantErr: true,
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
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"description":{"x":"y"}}`),
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
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
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
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
			},
			wantErr: true,
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
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"description":{"x":"y"}}`),
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
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
			},
			wantErr: true,
		},
		{
			name: "body empty",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(``),
			},
			wantErr: true,
		},
		{
			name: "no description in body",
			fields: fields{
				httpClient: new(mockHttp),
			},
			args: args{
				ctx: context.Background(),
				notifier: &domain.Notifier{
					Config: &domain.NotifierConfiguration{
						Secret: &domain.NotifierSecret{Token: "token"},
						Opts: map[string]interface{}{
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"summary":"abc"}`),
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
							"project_key": "abc", "issue_type": "xyz",
						},
					},
				},
				body: []byte(`{"summary":"abc","description":{"x":"y"}}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &jiraSimple{
				Client: &Client{HTTPClient: tt.fields.httpClient},
			}
			got, err := p.Send(tt.args.ctx, tt.args.notifier, tt.args.body)
			if tt.wantErr == false {
				if got.ID == "" {
					t.Errorf("jiraSimple.Send() ID missing in payload")
				}
				if got.Ok == false {
					t.Errorf("jiraSimple.Send() Ok == false")
				}
			}
			if err != nil != tt.wantErr {
				t.Errorf("jiraSimple.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpts_Extract(t *testing.T) {
	got := new(Opts)
	secret := &domain.NotifierSecret{Token: "token"}
	want := &Opts{
		ProjectKey: "abc",
		IssueType:  "xyz",
		Secret:     secret,
	}
	got.Extract(&domain.NotifierConfiguration{
		Opts: map[string]interface{}{
			"project_key": "abc",
			"issue_type":  "xyz",
		},
		Secret: secret,
	})
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Opts.Extract() = %v, want %v", got, want)
	}
}
