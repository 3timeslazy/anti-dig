package main

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/observability"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/http/server"
)

func Provide() *server.Server {
	var2 := config.NewConfig()
	var4_0 := observability.NewObservability(var2)

	var6_0 := flatten.NewListOfHandlers(var4_0.Metrics)

	var8_httpHandlers := []handlers.Handler{}
	var8_httpHandlers = append(var8_httpHandlers, var6_0.Handlers...)
	var9 := server.ServerParams{
		Config:		var2,
		Handlers:	var8_httpHandlers,
	}
	var1 := server.NewServer(var9)
	return var1
}
