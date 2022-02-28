package service

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/segmentio/ksuid"
)

type subscriptionService struct {
	subscriptionRepository domain.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepository domain.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		subscriptionRepository: subscriptionRepository,
	}
}

func (service *subscriptionService) Create(ctx context.Context, request *CreateSubscriptionRequest) (*domain.Subscription, error) {
	subscription := &domain.Subscription{
		ID:             ksuid.New().String(),
		EventType:      request.EventType,
		RuleExpression: request.RuleExpression,
		NotifierID:     request.NotifierID,
		Priority:       request.Priority,
		TemplateID:     request.TemplateID,
		TenantID:       request.TenantID,
		DateCreated:    time.Now(),
		DateUpdated:    time.Now(),
	}
	if err := service.subscriptionRepository.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}
