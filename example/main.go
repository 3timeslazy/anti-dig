package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
	"github.com/3timeslazy/anti-dig/example/cron"
	"github.com/3timeslazy/anti-dig/example/db"
	grpcserv "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	httpserv "github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/observability"
)

func main() {
	container := dig.New()

	err := container.Provide(observability.NewObservability)
	if err != nil {
		panic(err)
	}
	err = container.Provide(flatten.NewListOfHandlers)
	if err != nil {
		panic(err)
	}
	err = container.Provide(handlerv0.NewHandlerV0, dig.Group("http_handlers"))
	if err != nil {
		panic(err)
	}
	err = container.Provide(handlerv1.NewHandlerV1)
	if err != nil {
		panic(err)
	}
	err = container.Provide(cron.NewCron)
	if err != nil {
		panic(err)
	}
	err = container.Provide(db.NewDB)
	if err != nil {
		panic(err)
	}
	err = container.Provide(config.NewConfig)
	if err != nil {
		panic(err)
	}
	err = container.Provide(httpserv.NewServer)
	if err != nil {
		panic(err)
	}
	err = container.Provide(grpcserv.NewServer)
	if err != nil {
		panic(err)
	}
	err = container.Provide(consumer.New)
	if err != nil {
		panic(err)
	}
	err = container.Provide(queue.New1, dig.Name("queue_1"))
	if err != nil {
		panic(err)
	}
	err = container.Provide(queue.New2, dig.Name("queue_2"))
	if err != nil {
		panic(err)
	}
	err = container.Invoke(Run)
	if err != nil {
		panic(err)
	}
}

func Run(cron cron.Cron, httpsrv *httpserv.Server, grpcsrv *grpcserv.Server) {
	go cron.Start()
	go httpsrv.Serve()
	go grpcsrv.Serve()
}
