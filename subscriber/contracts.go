package subscriber

import "context"

type Service interface {
	Create(ctx context.Context, request *CreateRequest) (*Subscriber, error)
	GetByID(ctx context.Context, request *GetRequest) (*Subscriber, error)
}

type CreateRequest struct {
	Slug string `json:"slug"`
}

type GetRequest struct {
	ID string
}
