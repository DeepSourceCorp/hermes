package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	"os"

	"github.com/deepsourcelabs/hermes/config"
	handler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/service"
	sqlStore "github.com/deepsourcelabs/hermes/storage/sql"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartAsHTTPServer() error {
	log.Info("starting hermes in stateful mode...")
	db, err := gorm.Open(postgres.Open("postgres://hermes:password@localhost:5432/hermes"), &gorm.Config{SkipDefaultTransaction: false})
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
		messageHandler)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(":7272")
}

func main() {

	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	var isStateless = flag.Bool("stateless", true, "-stateless")
	var useEnv = flag.Bool("use-env", false, "-use-env")
	flag.Parse()

	cfg := new(config.AppConfig)
	if *useEnv {
		if err := cfg.ReadEnv(); err != nil {
			panic(err)
		}
	} else {
		if err := cfg.ReadYAML("./"); err != nil {
			panic(err)
		}
	}

	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	if *isStateless {
		log.Info("starting hermes in stateless mode...")
		if err := StartStatelessMode(cfg); err != nil {
			panic(err)
		}
		return
	}
	log.Info("starting hermes in stateful mode...")
	if err := StartAsHTTPServer(); err != nil {
		panic(err)
	}
}
