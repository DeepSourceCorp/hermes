package storage

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"gorm.io/gorm"
)

type templateStore struct {
	db *gorm.DB
}

func NewTemplateStore(db *gorm.DB) domain.TemplateRepository {
	return &templateStore{
		db: db,
	}
}

type Template struct {
	ID           string
	Pattern      string
	TemplateType string
	DateCreated  time.Time
	DateUpdated  time.Time
}

func newTemplate(t *domain.Template) *Template {
	return &Template{
		ID:           t.ID,
		Pattern:      t.Pattern,
		TemplateType: string(t.Type),
		DateCreated:  t.DateCreated,
		DateUpdated:  t.DateUpdated,
	}
}

func (t *Template) entity() *domain.Template {
	return &domain.Template{
		ID:      t.ID,
		Pattern: t.Pattern,
		Type:    domain.TemplateType(t.TemplateType),
	}
}

func (store *templateStore) Create(ctx context.Context, tmpl *domain.Template) domain.IError {
	t := newTemplate(tmpl)
	if err := store.db.Create(t).Error; err != nil {
		return errDBErr(err.Error())
	}
	return nil
}

func (store *templateStore) GetByID(ctx context.Context, id string) (*domain.Template, domain.IError) {
	tmpl := &Template{}
	if err := store.db.First(&tmpl, &Template{ID: id}).Error; err != nil {
		return nil, errDBErr(err.Error())
	}
	if tmpl != nil {
		return tmpl.entity(), nil
	}
	return nil, nil
}
