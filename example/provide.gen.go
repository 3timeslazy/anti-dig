package main

import (
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/cron"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/observability"
)

func Provide() (cron.Cron, *server.Server, *grpcserver.Server) {
	var2, err := db.NewDB()
	if err != nil {
		return nil, nil, nil
	}
	var4_queue1 := queue.New1()
	var5_queue2 := queue.New2()
	var6 := consumer.ConsumerParams{
		Queue1:	var4_queue1,
		Queue2:	var5_queue2,
	}
	var3 := consumer.New(var6)
	var7 := config.NewConfig()
	var1 := cron.NewCron(var2, var3, var7)
	var10_0 := observability.NewObservability(var7)

	var12_0 := flatten.NewListOfHandlers(var10_0.Metrics)

	var13_httpHandlers_1, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		return nil, nil, nil
	}
	var14_httpHandlers := []handlers.Handler{
		var13_httpHandlers_1,
	}
	var14_httpHandlers = append(var14_httpHandlers, var12_0.Handlers...)
	var15 := server.ServerParams{
		Config:		var7,
		Handlers:	var14_httpHandlers,
	}
	var8 := server.NewServer(var15)
	var17_0 := handlerv1.NewHandlerV1()

	var19_grpcHandlers := []handlers.Handler{
		var17_0.Handler,
	}
	var20 := grpcserver.ServerParams{
		Handlers: var19_grpcHandlers,
	}
	var16 := grpcserver.NewServer(var20)
	return var1, var8, var16
}
