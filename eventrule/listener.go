package eventrule

import (
	"context"
	"errors"

	"github.com/deepsourcelabs/hermes/event"
	"github.com/deepsourcelabs/hermes/rule"
	"github.com/deepsourcelabs/hermes/subscription"
)

type Listener interface {
	OnEvent()
}

type listener struct {
	svcSubscription subscription.Service
	svcRule         rule.Service
}

func (p *listener) OnEvent(ctx context.Context, event *event.Event) error {
	rules, errs := p.rulesForEvent(ctx, event)
	if len(errs) > 0 {
		return errors.New("error while loading rules")
	}

}

func (p *listener) rulesForEvent(ctx context.Context, event *event.Event) (rules []rule.SerializableRule, errors []error) {
	chRules := make(chan []rule.SerializableRule, 10)
	chErr := make(chan error, 10)

	request := subscription.GetAllRequest{
		SubscriberID: event.SubscriberID,
	}
	subscriptions, err := p.svcSubscription.GetAll(ctx, &request)
	if err != nil {
		return []rule.SerializableRule{}, nil
	}
	for _, v := range subscriptions {
		go p.rulesForSubscription(ctx, chRules, chErr, v.SubscriberID)
	}
	for range subscriptions {
		select {
		case r := <-chRules:
			rules = append(rules, r...)
		case err := <-chErr:
			errors = append(errors, err)
		}
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
