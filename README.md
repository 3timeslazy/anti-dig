# :carpentry_saw: anti-dig

**anti-dig** is a drop-in replacement for `go.uber.org/dig`, providing 100% compatibility. While it performs the same function as dig, it introduces a new approach. Instead of executing providers and passing them to other providers at runtime, it generates a file with the code calling your providers in the correct order.

# Usage

To use **anti-dig**, follow these steps:

1. Replace `go.uber.org/dig` with `github.com/3timeslazy/anti-dig`
```go
package main

import (
  // Replace "go.uber.org/dig" here
  dig "github.com/3timeslazy/anti-dig"
)

func main() {
  container := dig.New()
  
  // dig.Provide() calls here
  
  container.Invoke(run)
}
```

2. Run your code. It will generate a new file containing all the dependencies in a single file.

# Example
For an example illustrating the usage of anti-dig, refer to the code provided [here](https://github.com/3timeslazy/anti-dig/blob/main/example/main.go)

The original code:
```go
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
```

Will be transformed into:
```go
package main

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
	"github.com/3timeslazy/anti-dig/example/cron"
	"github.com/3timeslazy/anti-dig/example/db"
	grpcserver "github.com/3timeslazy/anti-dig/example/grpc/server"
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	"github.com/3timeslazy/anti-dig/example/http/server"
	"github.com/3timeslazy/anti-dig/example/observability"
)

func Provide() (cron.Cron, *server.Server, *grpcserver.Server) {
	config := config.NewConfig()
	db, err := db.NewDB(config)
	if err != nil {
		return nil, nil, nil
	}
	queue1 := queue.New1()
	queue := queue.New2()
	consumerParams := consumer.ConsumerParams{
		Queue1: queue1,
		Queue2: queue,
	}
	consumer := consumer.New(consumerParams)
	cron := cron.NewCron(db, consumer, config)
	observability := observability.NewObservability(config)

	listOfHandlers := flatten.NewListOfHandlers(observability.Metrics)

	handler, err := handlerv0.NewHandlerV0(db)
	if err != nil {
		return nil, nil, nil
	}
	httpHandlers := []handlers.Handler{
		handler,
	}
	httpHandlers = append(httpHandlers, listOfHandlers.Handlers...)
	serverParams := server.ServerParams{
		Config:   config,
		Handlers: httpHandlers,
	}
	server := server.NewServer(serverParams)
	handlerV1 := handlerv1.NewHandlerV1()

	grpcHandlers := []handlers.Handler{
		handlerV1.Handler,
	}
	grpcserverServerParams := grpcserver.ServerParams{
		Handlers: grpcHandlers,
	}
	grpcserverServer := grpcserver.NewServer(grpcserverServerParams)
	return cron, server, grpcserverServer
}
```

In the generated code, all the dependencies will be clearly defined and ordered according to their relationships.

Before running the generated code, make sure to perform the following steps:
1. Fix any error handling in the generated file, if required
2. Replace the `Invoke(Run)` statement in the original file with `Run(Provide())`
3. Remove the go.uber.org/dig dependency from your project.

Following these steps will allow you to execute the generated code, providing a clear view of the dependencies between types.

It is recommended to rename variables and organize the code in a way that suits your project's conventions and readability. Although ChatGPT can assist with this task, I wouldn't recommend using it with proprietary code
