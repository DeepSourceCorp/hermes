package rule

import "context"

type Repository interface {
	Create(ctx context.Context, sRule *SerializableRule) (*SerializableRule, error)
	GetByID(ctx context.Context, id string) (*SerializableRule, error)
	GetAll(ctx context.Context, subscriberID string) ([]SerializableRule, error)
}
