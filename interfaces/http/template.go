package http

import (
	"net/http"

	"github.com/deepsourcelabs/hermes/service"
	"github.com/labstack/echo/v4"
)

type templateHandler struct {
	templateService service.TemplateService
}

func NewTemplateHandler(templateService service.TemplateService) TemplateHandler {
	return &templateHandler{
		templateService: templateService,
	}
}

func (handler *templateHandler) PostTemplate() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		request := new(service.CreateTemplateRequest)
		if err := c.Bind(request); err != nil {
			return c.JSON(http.StatusBadRequest, "")
		}
		response, err := handler.templateService.Create(ctx, request)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "")
		}
		return c.JSON(http.StatusOK, response)
	}
}
