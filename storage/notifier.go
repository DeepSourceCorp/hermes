package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/deepsourcelabs/hermes/domain"
	"gorm.io/gorm"
)

type notifierStore struct {
	db *gorm.DB
}

// NewNotifierStore returns the GORM storage for Notifier.
func NewNotifierStore(db *gorm.DB) domain.NotifierRepository {
	return &notifierStore{
		db: db,
	}
}

// Notifier is the GORM representation of the Notifier.  This is done to ensure GORM
// specific tags, etc do not leak into the domain model and also to do serde within
// the implementation.
type Notifier struct {
	// ID is usually a KSUID primary key
	ID string

	// Type is the string representation of a NotifierType.  Eg: slack,JIRA, etc.
	Type string

	// Secret is the string representation from the domain.Secret type.
	Secret string

	// Extra is a serialized representation of domain.Extra
	Extras string

	DateCreated time.Time
	DateUpdated time.Time
}

// newNotifier returns a GORM representation of the Notifier from the domain entity.
func newNotifier(n *domain.Notifier) (*Notifier, error) {
	secret, err := json.Marshal(n.Config.Secret)
	if err != nil {
		return nil, err
	}
	extra, err := json.Marshal(n.Config.Opts)
	if err != nil {
		return nil, err
	}

	return &Notifier{
		ID:     n.ID,
		Type:   string(n.Type),
		Secret: string(secret),
		Extras: string(extra),
	}, nil
}

// entity returns the domain entity for a Notifier from the GORM specific representation
// of a Notifier.
func (n *Notifier) entity() (*domain.Notifier, error) {
	// deserialize secret
	secret := new(domain.NotifierSecret)
	if err := json.Unmarshal([]byte(n.Secret), &secret); err != nil {
		return nil, err
	}

	// deserialize extras
	opts := map[string]interface{}{}
	if err := json.Unmarshal([]byte(n.Extras), &opts); err != nil {
		return nil, err
	}

	return &domain.Notifier{
		ID: n.ID,
		Config: &domain.NotifierConfiguration{
			Secret: secret,
			Opts:   opts,
		},
	}, nil
}

// -- CRUD --

// Create inserts a new Notifier for persistance.
func (store *notifierStore) Create(ctx context.Context, notifier *domain.Notifier) error {
	n, err := newNotifier(notifier)
	if err != nil {
		return err
	}
	return store.db.Create(n).Error
}

func (store *notifierStore) Update(ctx context.Context, notifier *domain.Notifier) error {
	// Todo: Implement this.
	return nil
}
func (store *notifierStore) Delete(ctx context.Context, id string) error {
	// Todo: Implement this.
	return nil
}
func (store *notifierStore) GetByID(ctx context.Context, tenantID, id string) (*domain.Notifier, error) {
	notifier := &Notifier{}
	if err := store.db.First(&notifier, Notifier{ID: id}).Error; err != nil {
		return nil, err
	}
	if notifier != nil {
		return notifier.entity()
	}
	return nil, nil
}

func (store *notifierStore) GetAll(ctx context.Context, tenantID string) ([]domain.Notifier, error) {
	return []domain.Notifier{}, nil
}
