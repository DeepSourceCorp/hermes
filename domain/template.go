package domain

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/hoisie/mustache"
)

type TemplateType string

const (
	TemplateTypeText       TemplateType = "text"
	TemplateTypeMustache   TemplateType = "mustache"
	TemplateTypeGoTemplate TemplateType = "gotmpl"
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
		return &mustacheTemplater{}
	case TemplateTypeGoTemplate:
		return &goTemplater{}
	}
	return nil
}

type mustacheTemplater struct{}
type goTemplater struct{}

func (*mustacheTemplater) Execute(pattern string, params interface{}) ([]byte, error) {
	str := mustache.Render(pattern, params)
	return []byte(str), nil
}

func (*goTemplater) Execute(pattern string, params interface{}) ([]byte, error) {
	tmpl, err := template.New("template").Parse(pattern)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, params)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
