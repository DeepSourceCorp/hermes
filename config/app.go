package config

import (
	"errors"
	"os"
	"path"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	// Server configuration
	Port        int    `mapstructure:"PORT" yaml:"PORT"`
	TemplateDir string `mapstructure:"TEMPLATE_DIR" yaml:"TEMPLATE_DIR"`
}

const TagName = "mapstructure"

func (config *AppConfig) ReadEnv() error {
	m, err := env2Map(*config)
	if err != nil {
		return err
	}
	return mapstructure.Decode(m, config)
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

func (config *AppConfig) ReadYAML(configPath string) error {
	configBytes, err := os.ReadFile(path.Join(configPath, "./config.yaml"))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, config)
}
