package service

import (
	"context"

	"github.com/deepsourcelabs/hermes/domain"
	log "github.com/sirupsen/logrus"
)

type MessageService interface {
	Send(ctx context.Context, request *SendMessageRequest) (domain.Messages, domain.IError)
}

type messageService struct {
	templateRepository domain.TemplateRepository
}

func NewMessageService(templateRepository domain.TemplateRepository) MessageService {
	return &messageService{
		templateRepository: templateRepository,
	}
}

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

func (service *messageService) Send(
	ctx context.Context,
	request *SendMessageRequest,
) (domain.Messages, domain.IError) {
	messages := []domain.Message{}
	for _, recipient := range request.Recipients {
		notifier, err := service.getNotifier(ctx, recipient.Notifier)
		if err != nil {
			log.Errorf("Failed to get notifier %v: %v", recipient.Notifier, err)
			return domain.Messages{}, err
		}

		template, err := service.getTemplate(ctx, recipient.Template)
		if err != nil {
			log.Errorf("Failed to get template %v: %v", recipient.Template, err)
			return domain.Messages{}, err
		}

		body, err := service.getBody(ctx, template, request.Payload)
		if err != nil {
			log.Errorf("Failed to get body for request: %v", err)
			return []domain.Message{}, err
		}
		provider := newProvider(recipient.Notifier.Type)

		message, err := provider.Send(ctx, notifier, body)
		if err != nil {
			log.Errorf("Failed to send message: %v", err)
			return nil, err
		}

		messages = append(messages, *message)
	}
	return messages, nil
}

func (*messageService) getNotifier(
	_ context.Context,
	n *domain.Notifier,
) (*domain.Notifier, domain.IError) {
	if n.ID == "" && n.Config == nil {
		return nil, errMandatoryParamsMissing("missing notifier")
	}
	return n, nil
}

func (service *messageService) getTemplate(
	ctx context.Context,
	t *domain.Template,
) (*domain.Template, domain.IError) {
	if service.templateRepository == nil && t.ID != "" {
		return nil, errStateless("templateRepository == nil")
	}

	if t.ID == "" && (t.Pattern == "" || t.Type == "") {
		return nil, errMandatoryParamsMissing("missing pattern")
	}

	if t.ID != "" {
		return service.templateRepository.GetByID(ctx, t.ID)
	}

	return &domain.Template{
		Type:    t.Type,
		Pattern: t.Pattern,
	}, nil
}

func (*messageService) getBody(
	_ context.Context, t *domain.Template, payload *map[string]interface{},
) ([]byte, domain.IError) {
	templater := t.GetTemplater()
	body, err := templater.Execute(t.Pattern, payload)
	if err != nil {
		log.Errorf("Failed to execute template: %v", err)
		return nil, errUnprocessable("template execution failed")
	}
	return body, nil
}
