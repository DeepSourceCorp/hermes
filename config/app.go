package config

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type PGConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"db"`
}

type AppConfig struct {
	// Server configuration
	Port        int       `yaml:"port"`
	TemplateDir string    `yaml:"templateDir"`
	Postgres    *PGConfig `yaml:"postgres"`
}

func (config *AppConfig) ReadYAML(configPath string) error {
	configBytes, err := os.ReadFile(path.Join(configPath, "./config.yaml"))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, config)
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
