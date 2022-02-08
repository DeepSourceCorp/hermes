package rule

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, request *CreateRequest) (*SerializableRule, error)
	GetByID(ctx context.Context, request *GetRequest) (*SerializableRule, error)
	Filter(ctx context.Context, request *FilterRequest) ([]SerializableRule, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (svc *service) Create(ctx context.Context, request *CreateRequest) (*SerializableRule, error) {
	return svc.repository.Create(
		ctx,
		&SerializableRule{
			SubscriberID:   request.SubscriberID,
			SubscriptionID: request.SubscriptionID,
			Trigger:        request.Trigger,
			Action:         request.Action,
		},
	)
}

func (svc *service) GetByID(ctx context.Context, request *GetRequest) (*SerializableRule, error) {
	return svc.repository.GetByID(ctx, request.ID)
}

func (svc *service) Filter(ctx context.Context, request *FilterRequest) ([]SerializableRule, error) {
	return svc.repository.GetAll(ctx, request.SubscriptionID)
}
