package main

import (
	"sync"

	"github.com/RichardKnop/machinery/v1/config"
	"github.com/deepsourcelabs/hermes/event"
	"github.com/deepsourcelabs/hermes/eventrule"

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

type App struct {
	http      *echo.Echo
	store     *redis.Client
	taskQueue *infrastructure.Machinery
}

func (hermes *App) InitTaskQueue(opts *config.Config) {
	machinery, err := infrastructure.GetMachineryServer(opts)
	if err != nil {
		panic(err)
	}
	hermes.taskQueue = machinery
}

func (hermes *App) InitStore(opts *redis.Options) {
	hermes.store = infrastructure.GetRedisClient(opts)
}

func (hermes *App) startWorkers() {
	hermes.taskQueue.StartWorker("rule-engine", 100)
}

func (hermes *App) startHTTPServer() {

	e := echo.New()
	hermes.http = e

	// e.Use(middleware.Logger())

	eventNotifier := event.NewNotifier(hermes.taskQueue)
	eventService := event.NewService(eventNotifier)
	eventHandler := httpHandler.NewEventHandler(eventService)

	subscriberStore := redisStore.NewSubscriberStore(hermes.store)
	subscriberService := subscriber.NewService(subscriberStore)
	subscriberHandler := httpHandler.NewSubscriberHandler(subscriberService)

	subscriptionStore := redisStore.NewSubscriptionStore(hermes.store)
	subscriptionService := subscription.NewService(subscriptionStore)
	subscriptionHandler := httpHandler.NewSubscriptionHandler(subscriptionService)

	ruleStore := redisStore.NewRuleStore(hermes.store)
	ruleService := rule.NewService(ruleStore)
	ruleHandler := httpHandler.NewRuleHandler(ruleService)

	eventRuleListener := eventrule.NewEventListener(subscriptionService, ruleService, hermes.taskQueue)
	eventRuleListener.RegisterListener()

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
	hermes := new(App)

	hermes.InitStore(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	hermes.InitTaskQueue(&config.Config{
		Broker:        "redis://localhost:6379",
		DefaultQueue:  "machinery_tasks",
		ResultBackend: "redis://localhost:6379",
	})

	// TODO: There is a nightmare at the moment.  SIGINT is only received by the first routine, so Cmd+C doesn't work.
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		defer wg.Done()
		hermes.startWorkers()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		hermes.startHTTPServer()
	}()

	wg.Wait()
}
