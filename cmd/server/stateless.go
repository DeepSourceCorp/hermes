package main

import (
	"fmt"

	"github.com/deepsourcelabs/hermes/config"
	handler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/service"
	configStore "github.com/deepsourcelabs/hermes/storage/config"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func StartStatelessMode(cfg *config.AppConfig, e *echo.Echo) error {
	if err := config.InitTemplateConfig(cfg.TemplateConfigPath); err != nil {
		log.Errorf("failed to intitialize configuration, err = %v", err)
		return err
	}
	go config.StartTemplateConfigWatcher(cfg.TemplateConfigPath)

	templateConfigFactory := config.NewTemplateConfigFactory()

	templateStore := configStore.NewTemplateStore(templateConfigFactory)

	messsageService := service.NewMessageService(templateStore)
	messageHandler := handler.NewMessageHandler(messsageService)

	providerService := service.NewProviderService()
	providerHandler := handler.NewProviderHandler(providerService)

	router := handler.NewStatelessRouter(messageHandler, providerHandler)

	router.AddRoutes(e)
	log.Info("starting hermes in stateless mode")
	return e.Start(fmt.Sprintf(":%d", cfg.Port))
}
