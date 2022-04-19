package domain

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/templater"
)

type TemplateType string

const (
	TemplateTypeText       TemplateType = "text"
	TemplateTypeMustache   TemplateType = "mustache"
	TemplateTypeGoTemplate TemplateType = "gotmpl"
	TemplateTypeJinja      TemplateType = "jinja"
)

type Template struct {
	ID                 string         `json:"id,omitempty"`
	Pattern            string         `json:"pattern,omitempty"`
	Type               TemplateType   `json:"type,omitempty"`
	DateCreated        time.Time      `json:"date_created,omitempty"`
	DateUpdated        time.Time      `json:"date_updated,omitempty"`
	SupportedProviders []ProviderType `json:"supported_providers"`
}

func (t *Template) IsSupported(provider ProviderType) bool {
	for _, v := range t.SupportedProviders {
		if v == provider {
			return true
		}
	}
	return false
}

type TemplateRepository interface {
	Create(ctx context.Context, template *Template) IError
	GetByID(ctx context.Context, id string) (*Template, IError)
}

type Templater interface {
	Execute(pattern string, params interface{}) ([]byte, error)
}

func (t *Template) GetTemplater() Templater {
	switch t.Type {
	case TemplateTypeMustache:
		return &templater.Mustache{}
	case TemplateTypeGoTemplate:
		return &templater.Go{}
	case TemplateTypeJinja:
		return &templater.Jinja{}
	}
	return nil
}
