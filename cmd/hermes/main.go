package main

import (
	"sync"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/deepsourcelabs/hermes/event"

	"github.com/deepsourcelabs/hermes/infrastructure"
	"github.com/deepsourcelabs/hermes/interfaces/http"
	httpHandler "github.com/deepsourcelabs/hermes/interfaces/http"
	"github.com/deepsourcelabs/hermes/rule"
	redisStore "github.com/deepsourcelabs/hermes/store/redis"
	"github.com/deepsourcelabs/hermes/subscriber"
	"github.com/deepsourcelabs/hermes/subscription"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func startWorkers() {
	machineryOpts := &config.Config{
		Broker:        "redis://localhost:6379",
		DefaultQueue:  "machinery_tasks",
		ResultBackend: "redis://localhost:6379",
	}
	machinery, err := infrastructure.GetMachineryServer(machineryOpts)
	if err != nil {
		panic(err)
	}
	machinery.StartWorker("rule-engine", 100)
}

func startHTTPServer() {
	redisOpts := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}

	redisClient := infrastructure.GetRedisClient(redisOpts)

	e := echo.New()

	// e.Use(middleware.Logger())

	subscriberStore := redisStore.NewSubscriberStore(redisClient)
	subscriptionStore := redisStore.NewSubscriptionStore(redisClient)
	ruleStore := redisStore.NewRuleStore(redisClient)

	subscriberService := subscriber.NewService(subscriberStore)
	subscriptionService := subscription.NewService(subscriptionStore)
	ruleService := rule.NewService(ruleStore)
	eventService := event.NewService(nil)

	subscriberHandler := httpHandler.NewSubscriberHandler(subscriberService)
	subscriptionHandler := httpHandler.NewSubscriptionHandler(subscriptionService)
	ruleHandler := httpHandler.NewRuleHandler(ruleService)
	eventHandler := httpHandler.NewEventHandler(eventService)

	subscriberRouter := http.NewRouter(
		subscriberHandler,
		subscriptionHandler,
		ruleHandler,
		eventHandler,
	)
	subscriberRouter.AddRoutes(e)
	e.Logger.Fatal(e.Start(":7272"))
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		startHTTPServer()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		startWorkers()
	}()

	wg.Wait()
}
