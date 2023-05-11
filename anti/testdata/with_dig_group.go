package main

import (
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
)

func Provide() *server.Server {
	var2 := config.NewConfig()
	var3, err := db.NewDB(var2)
	if err != nil {
		return nil
	}
	var4_httpHandlers_0, err := handlerv0.NewHandlerV0(var3)
	if err != nil {
		return nil
	}
	var5_httpHandlers := []handlers.Handler{
		var4_httpHandlers_0,
	}
	var6 := server.ServerParams{
		Config:		var2,
		Handlers:	var5_httpHandlers,
	}
	var1 := server.NewServer(var6)
	return var1
}
