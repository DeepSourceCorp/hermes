package service

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/segmentio/ksuid"
)

type CreateTemplateRequest struct {
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

type TemplateService interface {
	Create(ctx context.Context, request *CreateTemplateRequest) (*domain.Template, error)
}

type templateService struct {
	repository domain.TemplateRepository
}

func NewTemplateService(repository domain.TemplateRepository) TemplateService {
	return &templateService{
		repository: repository,
	}
}

func (svc *templateService) Create(ctx context.Context, request *CreateTemplateRequest) (*domain.Template, error) {
	template := &domain.Template{
		ID:          ksuid.New().String(),
		Pattern:     request.Pattern,
		Type:        domain.TemplateType(request.Type),
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
	if err := svc.repository.Create(ctx, template); err != nil {
		return nil, err
	}
	return template, nil
}
