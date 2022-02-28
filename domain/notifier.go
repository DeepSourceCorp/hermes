package domain

import (
	"context"
	"time"
)

type NotifierSecret struct {
	Token string `json:"token"`
}

type NotifierConfiguration struct {
	Secret *NotifierSecret        `json:"secret"`
	Opts   map[string]interface{} `json:"options"`
}

type Notifier struct {
	ID          string                 `json:"id,omitempty"`
	Config      *NotifierConfiguration `json:"config,omitempty"`
	TenantID    string                 `json:"tenant_id,omitempty"`
	DateCreated time.Time              `json:"date_created,omitempty"`
	DateUpdated time.Time              `json:"date_updated,omitempty"`
	Type        ProviderType           `json:"type"`
}

type NotifierRepository interface {
	Create(ctx context.Context, notifier *Notifier) error
	Update(ctx context.Context, notifier *Notifier) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, tenantID, id string) (*Notifier, error)
	GetAll(ctx context.Context, tenantID string) ([]Notifier, error)
}
