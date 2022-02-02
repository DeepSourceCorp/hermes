package rule

type CreateRequest struct {
	Trigger        Trigger `json:"trigger"`
	Action         Opts    `json:"action"`
	SubscriberID   string  `json:"subscriber_id"`
	SubscriptionID string  `json:"subscription_id"`
}

type GetRequest struct {
	SubscriberID   string `json:"subscriber_id"`
	SubscriptionID string `json:"subscription_id"`
	ID             string `json:"id"`
}
