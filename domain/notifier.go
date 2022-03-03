package domain

import (
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

type NotifierRepository interface{}
