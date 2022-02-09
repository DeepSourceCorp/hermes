package event

type Event struct {
	ID           string `json:"id"`
	ReceivedAt   int64  `json:"received_at"`
	EventType    string `json:"event_type"`
	SubscriberID string `json:"subscriber_id"`
	Payload      []byte `json:"payload"`
}

/*
--Event Arrives-->GetSubscriptions(OwnerID)-->GetRules(SubscriptionID)
--RuleEvaluator(CheckPayload against all events)
-- For MatchingRules-->TriggerAction(Event)
*/
