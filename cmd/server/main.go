package main

import (
	"flag"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/gommon/log"

	"os"

	"github.com/deepsourcelabs/hermes/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

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
	e.Use(middleware.Logger())
	e.HideBanner = true
	AddDefaultRoutes(cfg, e)
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)
	if *isStateless {
		if err := StartStatelessMode(cfg, e); err != nil {
			log.Error("failed to start hermes in stateless mode, exiting")
			os.Exit(1)
		}
		return
	}

	if err := StartStatefulMode(cfg, e); err != nil {
		log.Error("failed to start hermes in stateful mode, exiting")
		panic(err)
	}
}
