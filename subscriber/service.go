package subscriber

import (
	"context"
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
	return svc.repository.Create(
		ctx,
		&Subscriber{Slug: req.Slug},
	)
}

func (svc *service) GetByID(ctx context.Context, req *GetRequest) (*Subscriber, error) {
	return svc.repository.GetByID(ctx, req.ID)
}
