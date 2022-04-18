package config

import (
	"reflect"
	"testing"
)

func TestAppConfig_Validate(t *testing.T) {
	type fields struct {
		Port               int
		TemplateConfigPath string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "validate AppConfig",
			fields:  fields{Port: 7272, TemplateConfigPath: "./"},
			wantErr: false,
		},
		{
			name:    "validate AppConfig with PORT 0",
			fields:  fields{Port: 0, TemplateConfigPath: "./"},
			wantErr: true,
		},
		{
			name:    "validate AppConfig with TemplateConfigPath empty",
			fields:  fields{Port: 7272, TemplateConfigPath: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AppConfig{
				Port:               tt.fields.Port,
				TemplateConfigPath: tt.fields.TemplateConfigPath,
			}
			if err := config.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("AppConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGConfig_GetDSN(t *testing.T) {
	t.Run("get dsn", func(t *testing.T) {

		want := "postgres://hermes:password@localhost:5432/hermesDB"

		pgConfig := &PGConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "hermes",
			Password: "password",
			Database: "hermesDB",
		}

		if got := pgConfig.GetDSN(); got != want {
			t.Errorf("PGConfig.GetDSN() = %v, want %v", got, want)
		}
	})
}

func TestAppConfig_ReadEnv(t *testing.T) {
	t.Run("read env with valid env", func(t *testing.T) {
		t.Setenv("HERMES_port", "7272")
		t.Setenv("HERMES_template_config_path", "./")
		t.Setenv("HERMES_postgres__host", "localhost")
		t.Setenv("HERMES_postgres__port", "5432")
		t.Setenv("HERMES_postgres__user", "hermes")
		t.Setenv("HERMES_postgres__password", "password")
		t.Setenv("HERMES_postgres__db", "db")
		appConfig := AppConfig{}
		if err := appConfig.ReadEnv(); err != nil {
			t.Errorf("AppConfig.ReadEnv() unexpected error = %v", err)
		}

		if appConfig.Port != 7272 {
			t.Errorf("AppConfig.ReadEnv().Port = %v, want %v", appConfig.Port, 7272)
		}
		if appConfig.TemplateConfigPath != "./" {
			t.Errorf("AppConfig.ReadEnv().TemplateConfigPath = %v, want %v", appConfig.TemplateConfigPath, "./")
		}

		want := &PGConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "hermes",
			Password: "password",
			Database: "db",
		}

		if !reflect.DeepEqual(appConfig.Postgres, want) {
			t.Errorf("AppConfig.ReadEnv().Postgres = %v, want %v", appConfig.Postgres, want)
		}

	})
}
