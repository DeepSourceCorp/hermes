package domain

type Messages []Message

type Message struct {
	ID               string      `json:"id"`
	Payload          interface{} `json:"payload"`
	Ok               bool        `json:"ok"`
	ProviderResponse interface{} `json:"provider_response,omitempty"`
}
