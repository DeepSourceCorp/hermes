package http

import (
	"net/http"

	"github.com/deepsourcelabs/hermes/service"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type ProviderHandler interface {
	GetProviderHandler() echo.HandlerFunc
}

type providerHandler struct {
	providerService service.ProviderService
}

func NewProviderHandler(providerService service.ProviderService) ProviderHandler {
	return &providerHandler{
		providerService: providerService,
	}
}

func (handler *providerHandler) GetProviderHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.GetProviderReqeuest)
		if err := c.Bind(request); err != nil {
			log.Errorf("Failed to bind request while getting provider handler: %v", err)
			return c.JSON(http.StatusBadRequest, ErrorResponse{"request malformed"})
		}
		request.Token = c.Request().Header.Get("X-Notifier-Token")
		response, err := handler.providerService.GetProvider(ctx, request)
		if err != nil {
			return c.JSON(err.StatusCode(), ErrorResponse{err.Message()})
		}
		return c.JSON(http.StatusOK, response)
	}
}
