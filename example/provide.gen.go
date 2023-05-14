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

func Provide() (cron.Cron, *server.Server, *grpcserver.Server) {
	config := config.NewConfig()
	db, err := db.NewDB(config)
	if err != nil {
		return nil, nil, nil
	}
	queue1 := queue.New1()
	queue := queue.New2()
	var7 := consumer.ConsumerParams{
		Queue1: queue1,
		Queue2: queue,
	}
	consumer := consumer.New(var7)
	cron := cron.NewCron(db, consumer, config)
	observability := observability.NewObservability(config)

	listOfHandlers := flatten.NewListOfHandlers(observability.Metrics)

	handler, err := handlerv0.NewHandlerV0(db)
	if err != nil {
		return nil, nil, nil
	}
	var14_httpHandlers := []handlers.Handler{
		handler,
	}
	var14_httpHandlers = append(var14_httpHandlers, listOfHandlers.Handlers...)
	var15 := server.ServerParams{
		Config:   config,
		Handlers: var14_httpHandlers,
	}
	server := server.NewServer(var15)
	handlerV1 := handlerv1.NewHandlerV1()

	var19_grpcHandlers := []handlers.Handler{
		handlerV1.Handler,
	}
	var20 := grpcserver.ServerParams{
		Handlers: var19_grpcHandlers,
	}
	grpcserverServer := grpcserver.NewServer(var20)
	return cron, server, grpcserverServer
}
