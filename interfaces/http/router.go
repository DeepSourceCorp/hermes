package http

import (
	"github.com/labstack/echo/v4"
)

type EchoRouter interface {
	AddRoutes(*echo.Echo)
}

type router struct {
	templateHandler TemplateHandler
	messageHandler  MessageHandler
}

func NewRouter(
	templateHandler TemplateHandler,
	messagemessageHandler MessageHandler,
) EchoRouter {
	return &router{
		templateHandler: templateHandler,
		messageHandler:  messagemessageHandler,
	}
}

func (r *router) AddRoutes(e *echo.Echo) {
	// templates
	e.POST("/templates", r.templateHandler.PostTemplate())
	e.POST("/tenants/:tenant_id/messages", r.messageHandler.PostMessage())
}

type statelessRouter struct {
	messageHandler  MessageHandler
	providerHandler ProviderHandler
}

func NewStatelessRouter(
	messageHandler MessageHandler,
	providerHandler ProviderHandler,
) EchoRouter {
	return &statelessRouter{
		messageHandler:  messageHandler,
		providerHandler: providerHandler,
	}
}

func (r *statelessRouter) AddRoutes(e *echo.Echo) {
	e.POST("/messages", r.messageHandler.PostMessage())

	providers := e.Group("/providers")
	providers.GET("/:provider", r.providerHandler.GetProviderHandler())
}
