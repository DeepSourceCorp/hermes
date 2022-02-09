package subscription

import "context"

type Repository interface {
	Create(ctx context.Context, subscription *Subscription) (*Subscription, error)
	GetByID(ctx context.Context, id string) (*Subscription, error)
	GetAll(ctx context.Context, subscriberID string) ([]Subscription, error)
}
