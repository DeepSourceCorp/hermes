package webhook

import (
	"github.com/deepsourcelabs/hermes/templater"
)

type Webhook struct {
	Templater templater.ITemplater
}

const INTGR_TYPE_WEBHOOK = "webhook"

func (w *Webhook) Do(params interface{}) (interface{}, error) {
	return nil, nil
}
