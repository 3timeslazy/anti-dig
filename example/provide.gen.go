package main

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
	"github.com/3timeslazy/anti-dig/example/cron"
	"github.com/3timeslazy/anti-dig/example/db"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	"github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/observability"
)

func Provide() (cron.Cron, *server.Server, *grpcserver.Server, error) {
	config := config.NewConfig()
	db, err := db.NewDB(config)
	if err != nil {
		return nil, nil, nil, err
	}
	queue1 := queue.New1()
	queue := queue.New2()
	consumerParams := consumer.ConsumerParams{
		Queue1: queue1,
		Queue2: queue,
	}
	consumer := consumer.New(consumerParams)
	cron := cron.NewCron(db, consumer, config)
	observability := observability.NewObservability(config)

	listOfHandlers := flatten.NewListOfHandlers(observability.Metrics)

	handler, err := handlerv0.NewHandlerV0(db)
	if err != nil {
		return nil, nil, nil, err
	}
	httpHandlers := []handlers.Handler{
		handler,
	}
	httpHandlers = append(httpHandlers, listOfHandlers.Handlers...)
	serverParams := server.ServerParams{
		Config:   config,
		Handlers: httpHandlers,
	}
	server := server.NewServer(serverParams)
	handlerV1 := handlerv1.NewHandlerV1()

	grpcHandlers := []handlers.Handler{
		handlerV1.Handler,
	}
	grpcserverServerParams := grpcserver.ServerParams{
		Handlers: grpcHandlers,
	}
	grpcserverServer := grpcserver.NewServer(grpcserverServerParams)
	return cron, server, grpcserverServer, nil
}
