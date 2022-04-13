package service

import (
	"context"
	"testing"

	"github.com/deepsourcelabs/hermes/domain"
)

func TestProviderService(t *testing.T) {
	ps := NewProviderService()
	ctx := context.Background()

	type fields struct {
		token        string
		providerType domain.ProviderType
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "must work with jira",
			fields:  fields{token: "secret", providerType: domain.ProviderType("jira")},
			wantErr: false,
		},
		{
			name:    "must work with slack",
			fields:  fields{token: "secret", providerType: domain.ProviderType("slack")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &GetProviderRequest{
				Token: tt.fields.token,
				Type:  tt.fields.providerType,
			}
			if _, err := ps.GetProvider(ctx, req); (err != nil) != tt.wantErr {
				t.Errorf("SendMessageRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
