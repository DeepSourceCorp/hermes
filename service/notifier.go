package service

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/segmentio/ksuid"
)

type notifierService struct {
	notifierRepository domain.NotifierRepository
}

func NewNotifierService(notifierRepository domain.NotifierRepository) NotifierService {
	return &notifierService{
		notifierRepository: notifierRepository,
	}
}

func (service *notifierService) Create(ctx context.Context, request *CreateNotifierRequest) (*domain.Notifier, error) {
	secret := domain.NotifierSecret(*request.Configuration.Secret)
	opts := map[string]interface{}(request.Configuration.Opts)
	notifier := &domain.Notifier{
		ID:   ksuid.New().String(),
		Type: domain.ProviderType(request.Type),
		Config: &domain.NotifierConfiguration{
			Secret: &secret,
			Opts:   opts,
		},
		TenantID:    request.TenantID,
		DateCreated: time.Now(),
		DateUpdated: time.Now(),
	}
	if err := service.notifierRepository.Create(ctx, notifier); err != nil {
		return nil, err
	}

	return notifier, nil
}

func (service *notifierService) GetByID(ctx context.Context, request *GetNotifierRequest) (*domain.Notifier, error) {
	return service.notifierRepository.GetByID(ctx, request.TenantID, request.ID)
}
