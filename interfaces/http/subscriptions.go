package http

import (
	"errors"
	"net/http"

	service "github.com/deepsourcelabs/hermes/subscription"
	"github.com/labstack/echo/v4"
)

type SubscriptionHandler interface {
	PostSubscription() echo.HandlerFunc
	GetSubscription() echo.HandlerFunc
	FilterSubscriptions() echo.HandlerFunc
}

type subscriptionHandler struct {
	service.Service
}

func NewSubscriptionHandler(svc service.Service) SubscriptionHandler {
	return &subscriptionHandler{
		svc,
	}
}

func (handler *subscriptionHandler) PostSubscription() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.CreateRequest)
		if err := c.Bind(request); err != nil {
			return failBadRequest(c, err)
		}
		request.SubscriberID = c.Param("subscriberID")
		response, err := handler.Create(ctx, request)
		if err != nil {
			failInternal(c, err)
		}
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *subscriptionHandler) GetSubscription() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.GetRequest)
		request.SubscriberID = c.Param("subscriberID")
		request.ID = c.Param("id")
		if request.SubscriberID == "" || request.ID == "" {
			return failBadRequest(c, errors.New("mandatory params missing"))
		}

		response, err := handler.GetByID(ctx, request)
		if err != nil {
			return failInternal(c, err)
		}
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *subscriptionHandler) FilterSubscriptions() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.GetAllRequest)
		request.SubscriberID = c.Param("subscriberID")
		if request.SubscriberID == "" {
			return failBadRequest(c, errors.New("mandatory params missing"))
		}
		response, err := handler.GetAll(ctx, request)
		if err != nil {
			return failInternal(c, err)
		}
		return c.JSON(http.StatusOK, response)
	}
}
