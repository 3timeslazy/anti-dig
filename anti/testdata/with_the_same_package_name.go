package main

import (
	"github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers"
)

func Provide() (*server.Server, *grpcserver.Server) {
	var2 := config.NewConfig()
	var3, err := db.NewDB(var2)
	if err != nil {
		return nil, nil
	}
	var4_httpHandlers_0, err := handlerv0.NewHandlerV0(var3)
	if err != nil {
		return nil, nil
	}
	var5_httpHandlers := []handlers.Handler{
		var4_httpHandlers_0,
	}
	var6 := server.ServerParams{
		Config:		var2,
		Handlers:	var5_httpHandlers,
	}
	var1 := server.NewServer(var6)
	var8_0 := handlerv1.NewHandlerV1()

	var10_grpcHandlers := []handlers.Handler{
		var8_0.Handler,
	}
	var11 := grpcserver.ServerParams{
		Handlers: var10_grpcHandlers,
	}
	var7 := grpcserver.NewServer(var11)
	return var1, var7
}
