package main

import (
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/grpc/server"
)

func Provide() *server.Server {
	var2_0 := handlerv1.NewHandlerV1()

	var4_grpcHandlers := []handlers.Handler{
		var2_0.Handler,
	}
	var5 := server.ServerParams{
		Handlers: var4_grpcHandlers,
	}
	var1 := server.NewServer(var5)
	return var1
}
