package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

const (
	envPrefix = "HERMES_"
)

type PGConfig struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Database string `koanf:"db"`
}

func (pgConfig *PGConfig) GetDSN() string {
	// postgres://hermes:password@localhost:5432/hermes
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		pgConfig.User,
		pgConfig.Password,
		pgConfig.Host,
		pgConfig.Port,
		pgConfig.Database,
	)
}

type AppConfig struct {
	// Server configuration
	Port        int       `koanf:"port"`
	TemplateDir string    `koanf:"templatedir"`
	Postgres    *PGConfig `koanf:"postgres"`
}

func (config *AppConfig) ReadEnv() error {
	var k = koanf.New(".")
	k.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil)

	return k.Unmarshal("", config)
}

func (config *AppConfig) Validate() error {
	if config.Port == 0 {
		return errors.New("PORT not defined in env")
	}
	if config.TemplateDir == "" {
		return errors.New("TEMPLATE_CONFIG not defined in env")
	}
	return nil
}
