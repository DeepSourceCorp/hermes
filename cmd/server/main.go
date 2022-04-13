package main

import (
	"flag"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"os"

	"github.com/deepsourcelabs/hermes/config"
	"github.com/labstack/echo/v4"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(
		&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	)

	var isStateless = flag.Bool("stateless", true, "-stateless")

	flag.Parse()

	// Parse config
	cfg := new(config.AppConfig)
	if err := cfg.ReadEnv(); err != nil {
		log.Errorf("failed to initalize configuration, err=%v", err)
		os.Exit(1)
	}
	if err := cfg.Validate(); err != nil {
		log.Errorf("app configuration is invalid, err=%v", err)
		os.Exit(1)
	}

	// Initialize web server
	e := echo.New()
	e.HideBanner = true

	// Setup prometheus metrics handler
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	if *isStateless {
		if err := StartStatelessMode(cfg, e); err != nil {
			log.Error("failed to start hermes in stateless mode, exiting")
			os.Exit(1)
		}
		return
	}

	if err := StartStatefulMode(cfg, e); err != nil {
		panic(err)
	}
}
