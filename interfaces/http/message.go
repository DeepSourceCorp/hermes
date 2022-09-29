package http

import (
	"net/http"

	"github.com/deepsourcelabs/hermes/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type messageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) MessageHandler {
	return &messageHandler{
		messageService: messageService,
	}
}

func (handler *messageHandler) PostMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		var request = new(service.SendMessageRequest)
		if err := c.Bind(request); err != nil {
			log.Errorf("Failed to bind request while posting message: %v", err)
			return c.JSON(http.StatusBadRequest, ErrorResponse{"request malformed"})
		}
		if err := request.Validate(); err != nil {
			log.Errorf("Failed to validate request while posting message: %v", err)
			return c.JSON(err.StatusCode(), ErrorResponse{err.Message()})
		}
		response, err := handler.messageService.Send(ctx, request)
		if err != nil {
			log.Errorf("Failed to send request while posting message: %v", err)
			return c.JSON(err.StatusCode(), ErrorResponse{err.Message()})
		}
		return c.JSON(http.StatusOK, response)
	}
}
