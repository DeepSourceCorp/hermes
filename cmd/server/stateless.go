package main

import (
	"fmt"

	"github.com/deepsourcelabs/hermes/config"
	handler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/service"
	configStore "github.com/deepsourcelabs/hermes/storage/config"
	"github.com/labstack/echo/v4"
)

func StartStatelessMode(cfg *config.AppConfig) error {
	if err := config.InitTemplateConfig(cfg.TemplateDir); err != nil {
		return err
	}
	templateConfigFactory := config.NewTemplateConfigFactory()

	templateStore := configStore.NewTemplateStore(templateConfigFactory)

	messsageService := service.NewMessageService(templateStore)
	messageHandler := handler.NewMessageHandler(messsageService)

	providerService := service.NewProviderService()
	providerHandler := handler.NewProviderHandler(providerService)

	router := handler.NewStatelessRouter(messageHandler, providerHandler)

	e := echo.New()
	e.HideBanner = true
	router.AddRoutes(e)
	return e.Start(fmt.Sprintf(":%d", cfg.Port))
}
