package config

import (
	"os"
	"path"

	"github.com/deepsourcelabs/hermes/domain"
	"gopkg.in/yaml.v3"
)

var templateConfig *TemplateConfig

type Template struct {
	ID                 string                `mapstructure:"id,omitempty"`
	Path               string                `mapstructure:"path,omitempty"`
	Type               domain.TemplateType   `mapstructure:"type,omitempty"`
	SupportedProviders []domain.ProviderType `mapstructure:"supported_providers"`
}

type TemplateConfig struct {
	Templates []Template `mapstructure:"templates"`
}

func (config *TemplateConfig) ReadYAML(configPath string) error {
	configBytes, err := os.ReadFile(path.Join(configPath, "./template.yaml"))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, config)
}

func InitTemplateConfig(templateDir string) error {
	return templateConfig.ReadYAML(templateDir)
}

type TemplateConfigFactory interface {
	GetTemplateConfig() *TemplateConfig
}

type templateConfigFactory struct {
}

func NewTemplateConfigFactory() TemplateConfigFactory {
	return &templateConfigFactory{}
}

func (*templateConfigFactory) GetTemplateConfig() *TemplateConfig {
	return templateConfig
}
