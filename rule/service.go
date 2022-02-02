package rule

import (
	"context"
)

type Service interface {
	Create(ctx context.Context, request *CreateRequest) (*Rule, error)
	GetByID(ctx context.Context, request *GetRequest) (*Rule, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository,
	}
}

func (svc *service) Create(ctx context.Context, request *CreateRequest) (*Rule, error) {
	opts := &Opts{
		Type: request.Action.Type,
	}
	a := NewAction(opts)
	rule := &Rule{
		Trigger: request.Trigger,
		Action:  a,
	}

	if err := svc.repository.Create(ctx, rule); err != nil {
		return nil, err
	}

	return rule, nil
}

func (svc *service) GetByID(ctx context.Context, request *GetRequest) (*Rule, error) {
	rule, err := svc.repository.GetByID(ctx, request.SubscriberID, request.SubscriptionID, request.ID)
	if err != nil {
		return nil, err
	}
	return rule, nil
}
