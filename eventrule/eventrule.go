package eventrule

import (
	"github.com/deepsourcelabs/hermes/event"
	"github.com/deepsourcelabs/hermes/rule"
)

type EventRule struct {
	Event event.Event
	Rule  rule.Rule
}
