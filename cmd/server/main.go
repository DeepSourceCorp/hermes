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

	var isStateless = flag.Bool("stateless", true, "-stateless")

	flag.Parse()

	// Parse config
	cfg := new(config.AppConfig)
	if err := cfg.ReadEnv(); err != nil {
		panic(err)
	}
	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	// Initialize web server
	e := echo.New()
	e.HideBanner = true

	// Setup prometheus metrics handler
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	if *isStateless {
		if err := StartStatelessMode(cfg, e); err != nil {
			panic(err)
		}
		return
	}

	if err := StartStatefulMode(cfg, e); err != nil {
		panic(err)
	}
}
