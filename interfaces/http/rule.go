package http

import (
	"net/http"

	service "github.com/deepsourcelabs/hermes/rule"
	"github.com/labstack/echo/v4"
)

type RuleHandler interface {
	PostRule() echo.HandlerFunc
	GetRule() echo.HandlerFunc
}

type ruleHandler struct {
	service.Service
}

func NewRuleHandler(svc service.Service) RuleHandler {
	return &ruleHandler{
		svc,
	}
}

func (handler *ruleHandler) PostRule() echo.HandlerFunc {
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

func (handler *ruleHandler) GetRule() echo.HandlerFunc {
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
