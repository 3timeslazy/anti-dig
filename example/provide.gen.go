package main

import (
	"github.com/3timeslazy/anti-dig/example/handlers"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/cron"
	"github.com/3timeslazy/anti-dig/example/observability"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
	"github.com/3timeslazy/anti-dig/example/consumer"
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

	var14_0, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		return nil, nil, nil
	}

	var15_httpHandlers := []handlers.Handler{
		var14_0.Handler,
	}
	var15_httpHandlers = append(var15_httpHandlers, var12_0.Handlers...)
	var16 := server.ServerParams{
		Config:		var7,
		Handlers:	var15_httpHandlers,
	}
	var8 := server.NewServer(var16)
	var18_0 := handlerv1.NewHandlerV1()

	var20_grpcHandlers := []handlers.Handler{
		var18_0.Handler,
	}
	var21 := grpcserver.ServerParams{
		Handlers: var20_grpcHandlers,
	}
	var17 := grpcserver.NewServer(var21)
	return var1, var8, var17
}
