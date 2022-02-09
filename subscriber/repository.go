package subscriber

import "context"

type Repository interface {
	Create(ctx context.Context, subscriber *Subscriber) (*Subscriber, error)
	GetByID(ctx context.Context, id string) (*Subscriber, error)
}
