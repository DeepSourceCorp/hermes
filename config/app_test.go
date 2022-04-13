package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestGetDSN(t *testing.T) {
	pgConfig := &PGConfig{
		User:     "hermes",
		Password: "password",
		Host:     "localhost",
		Port:     5432,
		Database: "hermes",
	}
	got := pgConfig.GetDSN()
	want := "postgres://hermes:password@localhost:5432/hermes"

	if got != want {
		t.Errorf("dsn doesn't match, got: %s, want: %s\n", got, want)
	}
}

func TestReadEnv(t *testing.T) {
	pgConfig := &PGConfig{
		User:     "hermes",
		Password: "password",
		Host:     "localhost",
		Port:     5432,
		Database: "hermes",
	}

	conf := &AppConfig{
		Port:        8080,
		TemplateDir: "./templates",
		Postgres:    pgConfig,
	}

	err := conf.ReadEnv()
	assert.Nil(t, err)
}
