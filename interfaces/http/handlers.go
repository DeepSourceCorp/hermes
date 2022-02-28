package http

import "github.com/labstack/echo/v4"

type TenantHandler interface {
	PostTenant() echo.HandlerFunc
}

type TemplateHandler interface {
	PostTemplate() echo.HandlerFunc
}

type SubscriptionHandler interface {
	PostSubscription() echo.HandlerFunc
}

type NotiferHandler interface {
	PostNotifier() echo.HandlerFunc
	GetNotifier() echo.HandlerFunc
}

type MessageHandler interface {
	PostMessage() echo.HandlerFunc
}
