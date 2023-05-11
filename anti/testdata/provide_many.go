package main

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
)

func Provide() handlers.Handler {
	var3 := config.NewConfig()
	var2, err := db.NewDB(var3)
	if err != nil {
		return nil
	}
	var1, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		return nil
	}
	return var1
}
