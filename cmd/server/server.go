package main

import (
	"flag"
	"fmt"

	"github.com/deepsourcelabs/hermes/config"
	"github.com/deepsourcelabs/hermes/domain"
	handler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/service"
	cfgStore "github.com/deepsourcelabs/hermes/storage/config"
	sqlStore "github.com/deepsourcelabs/hermes/storage/sql"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartAsHTTPServer() error {
	db, err := gorm.Open(postgres.Open("postgres://hermes:password@localhost:5432/hermes"), &gorm.Config{})
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

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(":7272")
}

func StartInStatelessMode(cfg *config.AppConfig) error {
	templateConfig, err := loadTemplateConfig(cfg.TemplateConfigPath)
	if err != nil {
		return err
	}

	var templateStore domain.TemplateRepository

	if templateConfig != nil {
		templateStore = cfgStore.NewTemplateStore(templateConfig)
	}

	messsageService := service.NewMessageService(templateStore)
	messageHandler := handler.NewMessageHandler(messsageService)

	router := handler.NewStatelessRouter(messageHandler)

	e := echo.New()
	router.AddRoutes(e)
	return e.Start(fmt.Sprintf(":%d", cfg.Port))
}

func loadTemplateConfig(path string) (*config.TemplateCfg, error) {
	var templateConfig = new(config.TemplateCfg)

	viper.AddConfigPath(path)
	viper.SetConfigName("templates")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(templateConfig); err != nil {
		return nil, err
	}
	return templateConfig, nil
}

func loadAppConfig(path string) (*config.AppConfig, error) {
	var appConfig = new(config.AppConfig)

	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(appConfig); err != nil {
		return nil, err
	}
	return appConfig, nil
}

func main() {
	cfg, err := loadAppConfig(".")
	if err != nil {
		panic(err)
	}

	var isStateless = flag.Bool("stateless", true, "foobar")
	if *isStateless {
		if err := StartInStatelessMode(cfg); err != nil {
			panic(err)
		}
	}
	if err := StartAsHTTPServer(); err != nil {
		panic(err)
	}
}
