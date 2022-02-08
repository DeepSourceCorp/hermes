package http

import (
	"encoding/json"
	"net/http"

	"github.com/deepsourcelabs/hermes/event"
	"github.com/labstack/echo/v4"
)

type EventHandler interface {
	PostEvent() echo.HandlerFunc
}

type eventHandler struct {
	service event.Service
}

func NewEventHandler(svc event.Service) *eventHandler {
	return &eventHandler{
		service: svc,
	}
}

type Event struct {
	SubscriberID string          `json:"subscriber_id"`
	Type         string          `json:"type"`
	Payload      json.RawMessage `json:"payload"`
}

func (handler *eventHandler) PostEvent() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := new(event.CreateEventRequest)
		ctx := c.Request().Context()
		err := c.Bind(request)
		if err != nil {
			return failBadRequest(c, err)
		}
		if err := handler.service.Create(ctx, request); err != nil {
			return failInternal(c, err)
		}
		return c.JSON(http.StatusAccepted, map[string]string{"status": "accepted"})
	}
}
