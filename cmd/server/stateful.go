package main

import (
	"errors"
	"fmt"

	"github.com/deepsourcelabs/hermes/config"
	handler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/service"
	sqlStore "github.com/deepsourcelabs/hermes/storage/sql"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartStatefulMode(cfg *config.AppConfig, e *echo.Echo) error {
	if cfg.Postgres == nil {
		return errors.New("postgres configuration not set")
	}
	db, err := gorm.Open(
		postgres.Open(cfg.Postgres.GetDSN()),
		&gorm.Config{SkipDefaultTransaction: true},
	)
	if err != nil {
		return err
	}

	templateStore := sqlStore.NewTemplateStore(db)
	templateService := service.NewTemplateService(templateStore)
	templateHandler := handler.NewTemplateHandler(templateService)

	messsageService := service.NewMessageService(templateStore)
	messageHandler := handler.NewMessageHandler(messsageService)

	router := handler.NewRouter(
		templateHandler,
		messageHandler,
	)

	router.AddRoutes(e)
	log.Info("starting hermes in stateful mode...")
	return e.Start(fmt.Sprintf(":%d", cfg.Port))
}
