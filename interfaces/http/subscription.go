package http

import (
	"net/http"

	"github.com/deepsourcelabs/hermes/service"
	"github.com/labstack/echo/v4"
)

type subscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(subscriptionService service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

func (handler *subscriptionHandler) PostSubscription() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.CreateSubscriptionRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		response, err := handler.subscriptionService.Create(ctx, request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, response)
	}
}
