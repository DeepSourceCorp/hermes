package eventrule

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/deepsourcelabs/hermes/event"
	"github.com/deepsourcelabs/hermes/infrastructure"
	"github.com/deepsourcelabs/hermes/rule"
	"github.com/deepsourcelabs/hermes/subscription"
)

type Listener interface {
	OnEvent(ctx context.Context, event *event.Event) error
	RegisterListener()
}

type listener struct {
	svcSubscription subscription.Service
	svcRule         rule.Service
	taskQueue       *infrastructure.Machinery
}

func NewEventListener(
	svcSubscription subscription.Service,
	svcRule rule.Service,
	taskQueue *infrastructure.Machinery,
) Listener {
	return &listener{
		svcSubscription: svcSubscription,
		svcRule:         svcRule,
		taskQueue:       taskQueue,
	}
}

func (p *listener) RegisterListener() {
	p.taskQueue.RegisterTask("event-created", func(ctx context.Context, eventString string) error {
		event := new(event.Event)
		if err := json.Unmarshal([]byte(eventString), event); err != nil {
			return err
		}
		return p.OnEvent(ctx, event)
	})
}

func (p *listener) OnEvent(ctx context.Context, event *event.Event) error {
	rules, err := p.rulesForEvent(ctx, event)
	if err != nil {
		return err
	}
	_ = rules
	// Todo:
	// 1. Convert SerializableRule to Rule.
	// 2. Fan out -> evaluate rule action -> trigger action if passing
	return nil
}

func (p *listener) rulesForEvent(ctx context.Context, event *event.Event) (rules []rule.SerializableRule, err error) {
	chRules := make(chan []rule.SerializableRule, 10)
	chErr := make(chan error, 10)

	errs := []error{}

	request := subscription.GetAllRequest{
		SubscriberID: event.SubscriberID,
	}
	subscriptions, err := p.svcSubscription.GetAll(ctx, &request)

	if err != nil {
		return []rule.SerializableRule{}, err
	}
	for _, v := range subscriptions {
		go p.rulesForSubscription(ctx, chRules, chErr, v.ID)
	}
	for range subscriptions {
		select {
		case r := <-chRules:
			rules = append(rules, r...)
		case err := <-chErr:
			errs = append(errs, err)
		}
	}
	if len(errs) > 1 {
		return []rule.SerializableRule{}, errors.New("something went wrong while retrieving rules")
	}
	return rules, nil
}

func (p *listener) rulesForSubscription(ctx context.Context, chRules chan<- []rule.SerializableRule, chErr chan<- error, subscriptionID string) {
	rules, err := p.svcRule.Filter(ctx, &rule.FilterRequest{SubscriptionID: subscriptionID})
	if err != nil {
		chErr <- err
		return
	}
	chRules <- rules
}
