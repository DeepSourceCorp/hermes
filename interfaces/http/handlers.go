package http

import "github.com/labstack/echo/v4"

type TemplateHandler interface {
	PostTemplate() echo.HandlerFunc
}

type MessageHandler interface {
	PostMessage() echo.HandlerFunc
}
