package config

import (
	"testing"
)

func TestAppConfig_Validate(t *testing.T) {
	type fields struct {
		Port        int
		TemplateDir string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "validate AppConfig",
			fields:  fields{Port: 7272, TemplateDir: "./"},
			wantErr: false,
		},
		{
			name:    "validate AppConfig with PORT 0",
			fields:  fields{Port: 0, TemplateDir: "./"},
			wantErr: true,
		},
		{
			name:    "validate AppConfig with TemplateDir empty",
			fields:  fields{Port: 7272, TemplateDir: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AppConfig{
				Port:        tt.fields.Port,
				TemplateDir: tt.fields.TemplateDir,
			}
			if err := config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AppConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
