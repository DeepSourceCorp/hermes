package config

import "github.com/deepsourcelabs/hermes/domain"

type Template struct {
	ID                 string                `mapstructure:"id,omitempty"`
	Path               string                `mapstructure:"path,omitempty"`
	Type               domain.TemplateType   `mapstructure:"type,omitempty"`
	SupportedProviders []domain.ProviderType `mapstructure:"supported_providers"`
}

type TemplateCfg struct {
	Templates []Template `mapstructure:"templates"`
}
