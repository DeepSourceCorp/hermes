package domain

type Messages []Message

type Message struct {
	ID               string      `json:"id"`
	Body             string      `json:"body"`
	Ok               bool        `json:"ok"`
	ProviderResponse interface{} `json:"provider_response,omitempty"`
}
