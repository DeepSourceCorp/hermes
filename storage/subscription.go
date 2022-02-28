package storage

import (
	"context"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"gorm.io/gorm"
)

type subscriptionStore struct {
	db *gorm.DB
}

func NewSubscriptionStore(db *gorm.DB) domain.SubscriptionRepository {
	return &subscriptionStore{
		db: db,
	}
}

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

func newSubscription(m *domain.Subscription) *Subscription {
	s := Subscription(*m)
	return &s
}

func (s *Subscription) entity() *domain.Subscription {
	m := domain.Subscription(*s)
	return &m
}

func (store *subscriptionStore) Create(ctx context.Context, m *domain.Subscription) error {
	s := newSubscription(m)
	return store.db.Create(s).Error
}

func (store *subscriptionStore) FilterBy(ctx context.Context, tenantID, eventType, notifierID string) ([]domain.Subscription, error) {
	var results []Subscription

	filters := Subscription{
		TenantID:  tenantID,
		EventType: eventType,
	}
	if notifierID != "" {
		filters.NotifierID = notifierID
	}

	if err := store.db.Where(filters).Find(&results).Error; err != nil {
		return []domain.Subscription{}, err
	}

	subscriptions := []domain.Subscription{}
	for _, v := range results {
		s := v.entity()
		subscriptions = append(subscriptions, *s)
	}

	return subscriptions, nil
}
