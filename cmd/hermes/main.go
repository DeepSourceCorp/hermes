package main

import (
	"fmt"

	redisInfra "github.com/deepsourcelabs/hermes/infrastructure/redis"
	"github.com/deepsourcelabs/hermes/interfaces/http"
	httpHandler "github.com/deepsourcelabs/hermes/interfaces/http"
	redisStore "github.com/deepsourcelabs/hermes/store/redis"
	"github.com/deepsourcelabs/hermes/subscriber"
	"github.com/deepsourcelabs/hermes/subscription"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

func StartHTTPServer() {
	redisOpts := &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
	redisClient := redisInfra.GetRedisClient(redisOpts)

	e := echo.New()

	subscriberStore := redisStore.NewSubscriberStore(redisClient)
	subscriptionStore := redisStore.NewSubscriptionStore(redisClient)

	subscriberService := subscriber.NewService(subscriberStore)
	subscriptionService := subscription.NewService(subscriptionStore)

	subscriberHandler := httpHandler.NewSubscriberHandler(subscriberService)
	subscriptionHandler := httpHandler.NewSubscriptionHandler(subscriptionService)

	subscriberRouter := http.NewRouter(subscriberHandler, subscriptionHandler)
	subscriberRouter.AddRoutes(e)
	e.Logger.Fatal(e.Start(":7272"))

}

func main() {
	fmt.Println("Starting server...")
	StartHTTPServer()
}
