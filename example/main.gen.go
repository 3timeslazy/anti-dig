package main

import (
	handlerv0 "github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	server "github.com/3timeslazy/anti-dig/example/server"
	handlers "github.com/3timeslazy/anti-dig/example/handlers"
	db "github.com/3timeslazy/anti-dig/example/db"
	config "github.com/3timeslazy/anti-dig/example/config"
	cron "github.com/3timeslazy/anti-dig/example/cron"
	observability "github.com/3timeslazy/anti-dig/example/observability"
	flatten "github.com/3timeslazy/anti-dig/example/handlers/flatten"
	handlerv1 "github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
)

func main() {
	var2, err := db.NewDB()
	if err != nil {
		panic(err)
	}
	var3 := config.NewConfig()
	var1 := cron.NewCron(var2, var3)
	var6_0 := observability.NewObservability(var3)
	var5 := var6_0.Metrics
	var7 := var6_0.Profiler
	var8_0 := flatten.NewListOfHandlers(var5)
	var9_0 := var8_0.Handlers
	var10_0, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		panic(err)
	}
	var9_1 := var10_0.Handler
	var10_1 := handlerv1.NewHandlerV1()
	var9_2 := var10_1.Handler
	var11 := []handlers.Handler{
		var9_1,
		var9_2,
	}
	var11 = append(var11, var9_0...)
	var12 := server.ServerParams{
		Handlers: var11,
	}
	var4 := server.NewServer(var12)
	FIXME(var1, var4)
}

