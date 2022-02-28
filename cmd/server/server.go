package main

import (
	handler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/service"
	"github.com/deepsourcelabs/hermes/storage"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartAsHTTPServer() error {
	db, err := gorm.Open(postgres.Open("postgres://hermes:password@localhost:5432/hermes"), &gorm.Config{})
	if err != nil {
		return err
	}

	templateStore := storage.NewTemplateStore(db)
	templateService := service.NewTemplateService(templateStore)
	templateHandler := handler.NewTemplateHandler(templateService)

	subscriptionStore := storage.NewSubscriptionStore(db)
	subscriptionService := service.NewSubscriptionService(subscriptionStore)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	notifierStore := storage.NewNotifierStore(db)
	notiferService := service.NewNotifierService(notifierStore)
	notifierHandler := handler.NewNotiferHandler(notiferService)

	messsageService := service.NewMessageService(templateStore)
	messageHandler := handler.NewMessageHandler(messsageService)

	router := handler.NewRouter(
		templateHandler,
		subscriptionHandler,
		notifierHandler,
		messageHandler,
	)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(":7272")
}

func main() {
	if err := StartAsHTTPServer(); err != nil {
		panic(err)
	}
}
