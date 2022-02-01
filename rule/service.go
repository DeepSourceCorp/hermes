package rule

import "context"

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

	return nil, nil
}

func (svc *service) GetByID(ctx context.Context, request *GetRequest) (*Rule, error) {
	return nil, nil
}
