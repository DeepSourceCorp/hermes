package subscriber

import (
	"context"

	"github.com/segmentio/ksuid"
)

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (svc *service) Create(ctx context.Context, req *CreateRequest) (*Subscriber, error) {
	subscriber := &Subscriber{
		ID:   ksuid.New().String(),
		Slug: req.Slug,
	}

	if err := svc.repository.Create(ctx, subscriber); err != nil {
		return nil, err
	}

	return subscriber, nil
}

func (svc *service) GetByID(ctx context.Context, req *GetRequest) (*Subscriber, error) {
	subscriber, err := svc.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return subscriber, nil
}
