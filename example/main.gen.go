package main

import (
	config "github.com/3timeslazy/anti-dig/example/config"
	cron "github.com/3timeslazy/anti-dig/example/cron"
	db "github.com/3timeslazy/anti-dig/example/db"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	handlers "github.com/3timeslazy/anti-dig/example/handlers"
	flatten "github.com/3timeslazy/anti-dig/example/handlers/flatten"
	handlerv0 "github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	handlerv1 "github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	server "github.com/3timeslazy/anti-dig/example/http/server"
	observability "github.com/3timeslazy/anti-dig/example/observability"
)

func main() {
	var2, err := db.NewDB()
	if err != nil {
		panic(err)
	}
	var3 := config.NewConfig()
	var1 := cron.NewCron(var2, var3)
	var6_0 := observability.NewObservability(var3)
	var8_0 := flatten.NewListOfHandlers(var6_0.Metrics)
	var10_0, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		panic(err)
	}
	var10_1 := handlerv1.NewHandlerV1()
	var11 := []handlers.Handler{
		var10_0.Handler,
		var10_1.Handler,
	}
	var11 = append(var11, var8_0.Handlers...)
	var12 := server.ServerParams{
		Handlers: var11,
	}
	var4 := server.NewServer(var12)
	var14 := grpcserver.ServerParams{
		Handlers: var11,
	}
	var13 := grpcserver.NewServer(var14)
	FIXME(var1, var4, var13)
}
