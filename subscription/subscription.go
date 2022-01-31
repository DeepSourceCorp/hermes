package subscription

type Secret struct {
	Method string `json:"type"`
	Val    string `json:"val"`
}
type Subscription struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Secret       Secret `json:"secret"`
	SubscriberID string `json:"subscriber_id"`
}
