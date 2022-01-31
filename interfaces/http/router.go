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
}

func NewRouter(subscriberHandler SubscriberHandler, subscriptionHandler SubscriptionHandler) Router {
	return &router{
		subscriberHandler,
		subscriptionHandler,
	}
}

func (r *router) AddRoutes(e *echo.Echo) {
	e.POST("/subscribers", r.subscriberHandler.PostSubscriber())
	e.GET("/subscribers/:id", r.subscriberHandler.GetSubscriber())
	e.POST("/subscribers/:subscriberID/subscriptions", r.subscriptionHandler.PostSubscription())
	e.GET("/subscribers/:subscriberID/subscriptions/:id", r.subscriptionHandler.GetSubscription())
	e.GET("/subscribers/:subscriberID/subscriptions", r.subscriptionHandler.FilterSubscriptions())
}
