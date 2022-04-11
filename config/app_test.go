package config

import (
	"testing"
)

func TestAppConfig_ReadEnv(t *testing.T) {
	getEnv = func(key string) string {
		if key == "PORT" {
			return "7272"
		}
		return "./foo"
	}

	cfg := new(AppConfig)
	err := cfg.ReadEnv()
	if err != nil {
		t.Errorf("AppConfig.ReadEnv() error = %v, wantErr %v", err, false)
	}

	if cfg.Port != 7272 {
		t.Errorf("AppConfig.ReadEnv() AppConfig.Port = %v, want %v", cfg.Port, 7272)
	}

	if cfg.TemplateDir != "./foo" {
		t.Errorf("AppConfig.ReadEnv() AppConfig.TemplateDir = %v, want = %v", cfg.TemplateDir, "./foo")
	}

	getEnv = func(key string) string {
		return "foo"
	}

	err = cfg.ReadEnv()
	if err == nil {
		t.Errorf("AppConfig.ReadEnv() error = %v, wantErr %v", err, true)
	}

}

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
