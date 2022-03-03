package config

import "github.com/deepsourcelabs/hermes/domain"

type Template struct {
	ID                 string                `yaml:"id,omitempty"`
	Pattern            string                `yaml:"pattern,omitempty"`
	Type               domain.TemplateType   `yaml:"type,omitempty"`
	SupportedProviders []domain.ProviderType `yaml:"supported_providers"`
}

type Config struct {
	Templates []Template `yaml:"templates"`
}
