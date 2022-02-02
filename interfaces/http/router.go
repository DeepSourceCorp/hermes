package http

import (
	"github.com/labstack/echo/v4"
)

type Router interface {
	AddRoutes(*echo.Echo)
}

type router struct {
	subscriberHandler   SubscriberHandler
	subscriptionHandler SubscriptionHandler
	ruleHandler         RuleHandler
}

func NewRouter(subscriberHandler SubscriberHandler, subscriptionHandler SubscriptionHandler, ruleHandler RuleHandler) Router {
	return &router{
		subscriberHandler,
		subscriptionHandler,
		ruleHandler,
	}
}

func (r *router) AddRoutes(e *echo.Echo) {
	e.POST("/subscribers", r.subscriberHandler.PostSubscriber())
	e.GET("/subscribers/:id", r.subscriberHandler.GetSubscriber())
	e.POST("/subscribers/:subscriberID/subscriptions", r.subscriptionHandler.PostSubscription())
	e.GET("/subscribers/:subscriberID/subscriptions/:id", r.subscriptionHandler.GetSubscription())
	e.GET("/subscribers/:subscriberID/subscriptions", r.subscriptionHandler.FilterSubscriptions())
	e.GET("/subscribers/:subscriberID/subscriptions/:subscriptionID/rules/:id", r.ruleHandler.GetRule())
	e.POST("/rules", r.ruleHandler.PostRule())
}
