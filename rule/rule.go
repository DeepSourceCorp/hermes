package rule

type Rule struct {
	Trigger Trigger
	Action  Action
}

type SerializableRule struct {
	ID             string             `json:"id"`
	SubscriptionID string             `json:"subscription_id"`
	SubscriberID   string             `json:"subscriber_id"`
	Trigger        Trigger            `json:"trigger"`
	Action         SerializableAction `json:"action"`
}

func (sr *SerializableRule) Rule() *Rule {
	action := NewAction(&sr.Action)
	return &Rule{
		Trigger: sr.Trigger,
		Action:  action,
	}
}
