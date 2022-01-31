package http

import (
	"net/http"

	service "github.com/deepsourcelabs/hermes/subscriber"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type SubscriberHandler interface {
	PostSubscriber() echo.HandlerFunc
	GetSubscriber() echo.HandlerFunc
}

type subscriberHandler struct {
	service.Service
}

func NewSubscriberHandler(svc service.Service) SubscriberHandler {
	return &subscriberHandler{
		svc,
	}
}

const RESP_FAIL = "failed"

type ErrorResponse struct {
	Status string `json:"status"`
}

func (handler *subscriberHandler) PostSubscriber() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.CreateRequest)
		if err := c.Bind(request); err != nil {
			failBadRequest(c, err)
		}
		response, err := handler.Create(ctx, request)
		if err != nil {
			failInternal(c, err)
		}
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *subscriberHandler) GetSubscriber() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.GetRequest)

		request.ID = c.Param(("id"))

		response, err := handler.GetByID(ctx, request)
		if err != nil {
			return failInternal(c, err)
		}
		return c.JSON(http.StatusOK, response)
	}
}

func failBadRequest(c echo.Context, err error) error {
	log.Errorf("failed to parse request, error=%v", err)
	return c.JSON(
		http.StatusBadRequest,
		ErrorResponse{
			Status: RESP_FAIL,
		},
	)
}

func failInternal(c echo.Context, err error) error {
	log.Errorf("something went wrong, error= %v", err)
	return c.JSON(
		http.StatusInternalServerError,
		ErrorResponse{
			Status: RESP_FAIL,
		},
	)
}
