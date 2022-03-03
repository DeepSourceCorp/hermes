package main

import (
	"flag"

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

	messsageService := service.NewMessageService(templateStore)
	messageHandler := handler.NewMessageHandler(messsageService)

	router := handler.NewRouter(
		templateHandler,
		messageHandler,
	)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(":7272")
}

func StartInStatelessMode() error {
	messsageService := service.NewMessageService(nil)
	messageHandler := handler.NewMessageHandler(messsageService)

	router := handler.NewStatelessRouter(messageHandler)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(":7272")
}

func main() {
	var isStateless = flag.Bool("stateless", true, "foobar")
	if *isStateless {
		if err := StartInStatelessMode(); err != nil {
			panic(err)
		}
	}
	if err := StartAsHTTPServer(); err != nil {
		panic(err)
	}
}
