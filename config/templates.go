package config

import (
	"fmt"
	"os"
	"path"

	"github.com/deepsourcelabs/hermes/domain"
	"gopkg.in/yaml.v3"
)

var templateConfig *TemplateConfig

type Template struct {
	ID                 string                `yaml:"id,omitempty"`
	Path               string                `yaml:"path,omitempty"`
	Type               domain.TemplateType   `yaml:"type,omitempty"`
	SupportedProviders []domain.ProviderType `yaml:"supported_providers"`
}

type TemplateConfig struct {
	Templates []Template `yaml:"templates"`
}

func (tc *TemplateConfig) Validate() error {
	for _, t := range tc.Templates {
		if _, err := os.Stat(t.Path); err != nil {
			return fmt.Errorf("template %s not found at %s", t.ID, t.Path)
		}
	}
	return nil
}

func (config *TemplateConfig) ReadYAML(configPath string) error {
	configBytes, err := os.ReadFile(path.Join(configPath, "./template.yaml"))
	if err != nil {
		return err
	}
	return yaml.Unmarshal(configBytes, &config)
}

func InitTemplateConfig(templateConfigPath string) error {
	templateConfig = new(TemplateConfig)
	if err := templateConfig.ReadYAML(templateConfigPath); err != nil {
		return err
	}
	return templateConfig.Validate()
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
