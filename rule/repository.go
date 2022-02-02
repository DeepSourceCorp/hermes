package rule

import "context"

type Repository interface {
	Create(ctx context.Context, rule *Rule) error
	GetByID(ctx context.Context, subscriberID, subscriptionID, id string) (*Rule, error)
}
