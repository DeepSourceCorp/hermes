package subscription

import "context"

type CreateRequest struct {
	Type         string `json:"type"`
	Secret       Secret `json:"secret"`
	SubscriberID string `json:"subscriber_id"`
}

type GetRequest struct {
	SubscriberID string
	ID           string
}

type GetAllRequest struct {
	SubscriberID string
}

type Service interface {
	Create(ctx context.Context, request *CreateRequest) (*Subscription, error)
	GetByID(ctx context.Context, request *GetRequest) (*Subscription, error)
	GetAll(ctx context.Context, request *GetAllRequest) ([]Subscription, error)
}
