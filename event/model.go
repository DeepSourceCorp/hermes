package event

type Event struct {
	ID           string
	ReceivedAt   int64
	EventType    string
	SubscriberID string
	Payload      []byte
}

/*
--Event Arrives-->GetSubscriptions(OwnerID)-->GetRules(SubscriptionID)
--RuleEvaluator(CheckPayload against all events)
-- For MatchingRules-->TriggerAction(Event)
*/
