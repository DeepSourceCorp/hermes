package subscription

import (
	"context"

	"github.com/segmentio/ksuid"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (svc *service) Create(ctx context.Context, request *CreateRequest) (*Subscription, error) {
	subscription := &Subscription{
		Secret:       request.Secret,
		ID:           ksuid.New().String(),
		Type:         request.Type,
		SubscriberID: request.SubscriberID,
	}
	return svc.repository.Create(ctx, subscription)
}

func (svc *service) GetByID(ctx context.Context, request *GetRequest) (*Subscription, error) {
	return svc.repository.GetByID(ctx, request.SubscriberID, request.ID)
}

func (svc *service) GetAll(ctx context.Context, request *GetAllRequest) ([]Subscription, error) {
	return svc.repository.GetAll(ctx, request.SubscriberID)
}
