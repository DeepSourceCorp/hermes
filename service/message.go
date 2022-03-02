package service

import (
	"context"
	"net/http"

	"github.com/deepsourcelabs/hermes/domain"
	"github.com/deepsourcelabs/hermes/provider"
	"github.com/deepsourcelabs/hermes/provider/slack"
)

type messageService struct {
	templateRepository domain.TemplateRepository
}

func NewMessageService(templateRepository domain.TemplateRepository) MessageService {
	return &messageService{
		templateRepository: templateRepository,
	}
}

func (service *messageService) Send(
	ctx context.Context,
	request *SendMessageRequest,
) (domain.Messages, domain.IError) {
	messages := []domain.Message{}
	for _, recipient := range request.Recipients {
		notifier, err := service.getNotifier(ctx, recipient.Notifier)
		if err != nil {
			return domain.Messages{}, err
		}

		template, err := service.getTemplate(ctx, recipient.Template)
		if err != nil {
			return domain.Messages{}, err
		}

		body, err := service.getBody(ctx, template, request.Payload)
		if err != nil {
			return []domain.Message{}, err
		}
		provider := newProvider(recipient.Notifier.Type)

		message, err := provider.Send(ctx, notifier, body)
		if err != nil {
			return nil, err
		}

		messages = append(messages, *message)
	}
	return messages, nil
}

func (service *messageService) getNotifier(
	ctx context.Context,
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

func (service *messageService) getBody(
	ctx context.Context,
	t *domain.Template,
	payload *map[string]interface{},
) ([]byte, domain.IError) {
	templater := t.GetTemplater()
	body, err := templater.Execute(t.Pattern, payload)
	if err != nil {
		return nil, errUnprocessable("template execution failed")
	}
	return body, nil
}

func newProvider(providerType domain.ProviderType) provider.Provider {
	switch providerType {
	case domain.ProviderTypeSlack:
		return slack.NewSlackProvider(http.DefaultClient)
	}
	return nil
}
