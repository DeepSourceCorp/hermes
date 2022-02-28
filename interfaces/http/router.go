package http

import "github.com/labstack/echo/v4"

type EchoRouter interface {
	AddRoutes(*echo.Echo)
}

type router struct {
	templateHandler      TemplateHandler
	sunbscriptionHandler SubscriptionHandler
	notifierHandler      NotiferHandler
	messageHandler       MessageHandler
}

func NewRouter(
	templateHandler TemplateHandler,
	subscriptionHandler SubscriptionHandler,
	notifierHandler NotiferHandler,
	messagemessageHandler MessageHandler,
) EchoRouter {
	return &router{
		templateHandler:      templateHandler,
		sunbscriptionHandler: subscriptionHandler,
		notifierHandler:      notifierHandler,
		messageHandler:       messagemessageHandler,
	}
}

func (r *router) AddRoutes(e *echo.Echo) {

	//templates
	e.POST("/templates", r.templateHandler.PostTemplate())

	//subscriptions
	e.POST("/tenants/:tenant_id/subscriptions", r.sunbscriptionHandler.PostSubscription())

	//notifiers
	e.POST("/tenants/:tenant_id/notifiers", r.notifierHandler.PostNotifier())
	e.GET("tenants/:tenant_id/notifiers/:id", r.notifierHandler.GetNotifier())

	e.POST("/tenants/:tenant_id/messages", r.messageHandler.PostMessage())
}
