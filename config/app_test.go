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
	t.Setenv("HERMES_PORT", "7272")
	t.Setenv("HERMES_TEMPLATEDIR", "./")
	t.Setenv("HERMES_POSTGRES_PORT", "5432")
	t.Setenv("HERMES_POSTGRES_HOST", "localhost")
	t.Setenv("HERMES_POSTGRES_USER", "hermes")
	t.Setenv("HERMES_POSTGRES_PASSWORD", "password")
	t.Setenv("HERMES_POSTGRES_DB", "hermes")

	pgConfig := &PGConfig{
		User:     "hermes",
		Password: "password",
		Host:     "localhost",
		Port:     5432,
		Database: "hermes",
	}

	want := &AppConfig{
		Port:        7272,
		TemplateDir: "./",
		Postgres:    pgConfig,
	}

	conf := &AppConfig{}
	err := conf.ReadEnv()
	assert.Equal(t, want, conf)
	assert.Nil(t, err)
}
