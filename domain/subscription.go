package domain

import (
	"context"
	"time"
)

type Subscription struct {
	ID             string
	EventType      string
	RuleExpression string
	Priority       uint
	NotifierID     string
	TemplateID     string
	TenantID       string
	DateCreated    time.Time
	DateUpdated    time.Time
}

type SubscriptionRepository interface {
	Create(ctx context.Context, s *Subscription) error
	FilterBy(ctx context.Context, tenantID, eventType, notifierID string) ([]Subscription, error)
}
