package event

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateEventRequest) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (svc *service) Create(ctx context.Context, request *CreateEventRequest) error {
	return nil
}
