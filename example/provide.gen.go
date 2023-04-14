package main

import (
	"github.com/3timeslazy/anti-dig/example/config"
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
	var2, err := db.NewDB()
	if err != nil {
		return nil, nil, nil
	}
	var3 := config.NewConfig()
	var1 := cron.NewCron(var2, var3)
	var6_0 := observability.NewObservability(var3)

	var8_0 := flatten.NewListOfHandlers(var6_0.Metrics)

	var10_0, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		return nil, nil, nil
	}

	var11_httpHandlers := []handlers.Handler{
		var10_0.Handler,
	}
	var11_httpHandlers = append(var11_httpHandlers, var8_0.Handlers...)
	var12 := server.ServerParams{
		Config:   var3,
		Handlers: var11_httpHandlers,
	}
	var4 := server.NewServer(var12)
	var14_0 := handlerv1.NewHandlerV1()

	var16_grpcHandlers := []handlers.Handler{
		var14_0.Handler,
	}
	var17 := grpcserver.ServerParams{
		Handlers: var16_grpcHandlers,
	}
	var13 := grpcserver.NewServer(var17)
	return var1, var4, var13
}
