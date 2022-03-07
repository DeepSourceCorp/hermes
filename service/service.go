package service

import (
	"context"

	"github.com/deepsourcelabs/hermes/domain"
)

type CreateTemplateRequest struct {
	Pattern string `json:"pattern"`
	Type    string `json:"type"`
}

type TemplateService interface {
	Create(ctx context.Context, request *CreateTemplateRequest) (*domain.Template, error)
}

// Messages

type SendMessageRequest struct {
	TenantID   string                  `param:"tenant_id"`
	Payload    *map[string]interface{} `json:"payload"`
	Recipients []struct {
		Notifier *domain.Notifier `json:"notifier"`
		Template *domain.Template `json:"template"`
	} `json:"recipients"`
}

func (r *SendMessageRequest) Validate() domain.IError {
	if r.Payload == nil {
		return errMandatoryParamsMissing("empty payload")
	}
	if len(r.Recipients) < 1 {
		return errMinOneRecipient("no recipients defined")
	}
	for _, v := range r.Recipients {
		if v.Notifier == nil || v.Template == nil {
			return errRecipientMalformed("some recipients are not valid")
		}
	}
	return nil
}

type MessageService interface {
	Send(ctx context.Context, request *SendMessageRequest) (domain.Messages, domain.IError)
}
