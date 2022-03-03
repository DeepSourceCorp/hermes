package config

import (
	"context"
	"os"

	"github.com/deepsourcelabs/hermes/config"
	"github.com/deepsourcelabs/hermes/domain"
)

type templateStore struct {
	cfg        *config.TemplateCfg
	fnReadFile func(filename string) ([]byte, error)
}

func NewTemplateStore(cfg *config.TemplateCfg) domain.TemplateRepository {
	return &templateStore{
		cfg:        cfg,
		fnReadFile: os.ReadFile,
	}
}

// Create creates a new store.  WARNING: Config store does not support writes at the moment.
func (store *templateStore) Create(ctx context.Context, tmpl *domain.Template) domain.IError {
	return errDBErr("config store does not support write")
}

//
func (store *templateStore) GetByID(ctx context.Context, id string) (*domain.Template, domain.IError) {
	var t = new(config.Template)

	for _, v := range store.cfg.Templates {
		if v.ID == id {
			t = &v
			break
		}
	}

	pattern, err := store.fnReadFile(t.Path)
	if err != nil {
		return nil, errDBErr(err.Error())
	}
	return &domain.Template{
		ID:                 id,
		Pattern:            string(pattern),
		SupportedProviders: t.SupportedProviders,
		Type:               t.Type,
	}, nil
}
