package config

import (
	"context"
	"os"

	"github.com/deepsourcelabs/hermes/config"
	"github.com/deepsourcelabs/hermes/domain"

	log "github.com/sirupsen/logrus"
)

type templateStore struct {
	cfg        *config.TemplateConfig
	fnReadFile func(filename string) ([]byte, error)
}

func NewTemplateStore(templateConfigFactory config.TemplateConfigFactory) domain.TemplateRepository {
	return &templateStore{
		cfg:        templateConfigFactory.GetTemplateConfig(),
		fnReadFile: os.ReadFile,
	}
}

// Create creates a new store.  WARNING: Config store does not support writes at the moment.
func (store *templateStore) Create(ctx context.Context, tmpl *domain.Template) domain.IError {
	return errDBErr("config store does not support write")
}

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
		log.Errorf("Failed to read file %v: %v", t.Path, err)
		return nil, errDBErr(err.Error())
	}
	return &domain.Template{
		ID:                 id,
		Pattern:            string(pattern),
		SupportedProviders: t.SupportedProviders,
		Type:               t.Type,
	}, nil
}
