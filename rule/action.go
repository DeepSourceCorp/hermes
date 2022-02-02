package rule

import "github.com/deepsourcelabs/hermes/integrations/webhook"

type Action interface {
	Do(params interface{}) (results interface{}, err error)
	Name() string
	TemplateID() string
}

type Opts struct {
	Type       string `json:"type"`
	TemplateID string `json:"template_id"`
}

func NewAction(opts *Opts) Action {
	switch opts.Type {
	case webhook.INTGR_TYPE_WEBHOOK:
		return &webhook.Webhook{}
	}
	return nil
}
