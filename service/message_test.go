package service

import (
	"testing"

	"github.com/deepsourcelabs/hermes/domain"
)

func TestMessageService_Validate(t *testing.T) {
	noNotifier := []struct {
		Notifier *domain.Notifier `json:"notifier"`
		Template *domain.Template `json:"template"`
	}{
		{nil, &domain.Template{ID: "1", Pattern: "example_pattern", Type: "gotmpl"}},
	}

	noTemplate := []struct {
		Notifier *domain.Notifier `json:"notifier"`
		Template *domain.Template `json:"template"`
	}{
		{&domain.Notifier{}, nil},
	}

	validRecipients := []struct {
		Notifier *domain.Notifier `json:"notifier"`
		Template *domain.Template `json:"template"`
	}{
		{&domain.Notifier{Config: &domain.NotifierConfiguration{Secret: &domain.NotifierSecret{Token: "secret"}, Opts: map[string]interface{}{"channel": "general"}}}, &domain.Template{ID: "template_1", Pattern: "example_pattern", Type: "gotmpl"}},
	}

	type fields struct {
		TenantID   string
		Payload    *map[string]interface{}
		Recipients []struct {
			Notifier *domain.Notifier `json:"notifier"`
			Template *domain.Template `json:"template"`
		} `json:"recipients"`
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "must not work with empty fields",
			fields:  fields{},
			wantErr: true,
		},
		{
			name:    "must not work with empty payload",
			fields:  fields{TenantID: "1", Payload: nil, Recipients: nil},
			wantErr: true,
		},
		{
			name:    "must not work with empty recipients",
			fields:  fields{TenantID: "1", Payload: &map[string]interface{}{"example": "payload"}, Recipients: nil},
			wantErr: true,
		},
		{
			name:    "must not work with recipients with empty notifier",
			fields:  fields{TenantID: "1", Payload: &map[string]interface{}{"example": "payload"}, Recipients: noNotifier},
			wantErr: true,
		},
		{
			name:    "must not work with recipients with empty template",
			fields:  fields{TenantID: "1", Payload: &map[string]interface{}{"example": "payload"}, Recipients: noTemplate},
			wantErr: true,
		},
		{
			name:    "must work with valid request",
			fields:  fields{TenantID: "1", Payload: &map[string]interface{}{"example": "payload"}, Recipients: validRecipients},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &SendMessageRequest{
				TenantID:   tt.fields.TenantID,
				Payload:    tt.fields.Payload,
				Recipients: tt.fields.Recipients,
			}
			if err := req.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("SendMessageRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
