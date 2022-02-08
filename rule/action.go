package rule

import "github.com/deepsourcelabs/hermes/integrations/webhook"

type Action interface {
	Do(params interface{}) (results interface{}, err error)
}

type SerializableAction struct {
	Type       string `json:"type"`
	TemplateID string `json:"template_id"`
}

func NewAction(a *SerializableAction) Action {
	switch a.Type {
	case webhook.INTGR_TYPE_WEBHOOK:
		return &webhook.Webhook{}
	}
	return nil
}
