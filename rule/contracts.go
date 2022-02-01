package rule

type CreateRequest struct {
	Trigger        Trigger `json:"trigger"`
	Action         Opts
	SubscriberID   string
	SubscriptionID string
}

type GetRequest struct {
	SubscriberID   string `json:"subscriber_id"`
	SubscriptionID string `json:"subscription_id"`
	ID             string `json:"id"`
}
