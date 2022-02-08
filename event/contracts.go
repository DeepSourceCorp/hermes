package event

import "encoding/json"

type CreateEventRequest struct {
	SubscriberID string          `json:"subscriber_id"`
	Type         string          `json:"type"`
	Payload      json.RawMessage `json:"payload"`
}
