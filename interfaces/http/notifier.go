package http

import (
	"net/http"

	"github.com/deepsourcelabs/hermes/service"
	"github.com/labstack/echo/v4"
)

type notifierHandler struct {
	notifierService service.NotifierService
}

func NewNotiferHandler(notifierService service.NotifierService) NotiferHandler {
	return &notifierHandler{
		notifierService: notifierService,
	}
}

func (handler *notifierHandler) PostNotifier() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.CreateNotifierRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		response, err := handler.notifierService.Create(ctx, request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, response)
	}
}

func (handler *notifierHandler) GetNotifier() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.GetNotifierRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}
		response, err := handler.notifierService.GetByID(ctx, request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, response)
	}
}
