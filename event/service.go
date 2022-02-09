package event

import (
	"context"
	"time"

	// "github.com/deepsourcelabs/hermes/eventrule"

	"github.com/segmentio/ksuid"
)

type Service interface {
	Create(context.Context, *CreateEventRequest) error
}

type service struct {
	notifier Notifier
}

func NewService(notifier Notifier) Service {
	return &service{
		notifier: notifier,
	}
}

func (svc *service) Create(ctx context.Context, request *CreateEventRequest) error {
	event := &Event{
		ID:           ksuid.New().String(),
		ReceivedAt:   time.Now().Unix(),
		EventType:    request.Type,
		Payload:      request.Payload,
		SubscriberID: request.SubscriberID,
	}
	return svc.notifier.Dispatch(ctx, event)
}
