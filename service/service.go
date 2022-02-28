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

type CreateSubscriptionRequest struct {
	EventType      string `json:"event_type"`
	RuleExpression string `json:"rule_expression"`
	Priority       uint   `json:"priority"`
	NotifierID     string `json:"notifier_id"`
	TemplateID     string `json:"template_id"`
	TenantID       string `param:"tenant_id"` //fixme: echo bind tags should not ideally leak to service layer
}

type SubscriptionService interface {
	Create(ctx context.Context, request *CreateSubscriptionRequest) (*domain.Subscription, error)
}

type NotifierSecret struct {
	Token string `json:"token"`
}

type NotifierConfiguration struct {
	Type   domain.ProviderType    `json:"type"`
	Secret NotifierSecret         `json:"secret"`
	Opts   map[string]interface{} `json:"options"`
}

type CreateNotifierRequest struct {
	Configuration domain.NotifierConfiguration `json:"config"`
	TenantID      string                       `param:"tenant_id"`
	Type          domain.ProviderType          `json:"type"`
}

type GetNotifierRequest struct {
	TenantID string `param:"tenant_id"`
	ID       string `param:"id"`
}

type NotifierService interface {
	Create(ctx context.Context, request *CreateNotifierRequest) (*domain.Notifier, error)
	GetByID(ctx context.Context, request *GetNotifierRequest) (*domain.Notifier, error)
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
