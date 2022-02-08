package eventrule

import (
	"context"

	"github.com/deepsourcelabs/hermes/event"
	"github.com/deepsourcelabs/hermes/rule"
)

type EventRule struct {
	Event *event.Event
	Rule  *rule.Rule
}

func (er *EventRule) isPassing() bool {
	return false
}

func (er *EventRule) Execute(ctx context.Context) error {
	if er.isPassing() {
		er.Rule.Action.Do(er.Event.Payload)
	}
	return nil
}
