# :carpentry_saw: anti-dig

An anti-dependency-injection drop-in replacement toolkit for `go.uber.org/dig`.

## Why?

I have worked in many companies. In every one of them, I've seen someone using `go.uber.org/dig`. And in each case, after a while the team wanted to get rid of it, but it didn't always work because it always took a lot of time and effort. Faced with this problem again, I decided to write a tool to help others get rid of the library.

## Table of Contents

* [Guide](#guide)
* [Example](#example)
* [Configuration](#configuration)

## Guide

To use **anti-dig**, follow the steps:

1. Replace `go.uber.org/dig` with `github.com/3timeslazy/anti-dig` everywhere in your code


```go
import (
    dig "github.com/3timeslazy/anti-dig" // instead of "go.uber.org/dig"
)
```

2. Download the package by running `go mod tidy`

3. Run your `main()` function. It will generate a new file with explicit initialization of your dependencies

## Example

Suppose we have a `main()` function with several providers passed to dig
```go
package main

import (
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

 	dig "github.com/3timeslazy/anti-dig" // "go.uber.org/dig"
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

The code above will generate the following:
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

Before running the generated code, make sure to perform the following steps:
1. Fix any error handling in the generated file, if required
2. Replace the `Invoke(Run)` statement in the original file with `Run(Provide())`
3. Remove go.uber.org/dig from your project 🥳🥳🥳 

Following these steps will allow you to execute the generated code, providing a clear view of your dependencies

## Configuration

**anti-dig** has two configuration options.

### Optimise

This option is enabled by default. 

If true, then anti-dig will generate code with some 'optimisations' such as human-readable variables, deleted unused variables, etc.

```go
package main

import dig "github.com/3timeslazy/anti-dig"

func main() {
	dig.Anti.Optimise(true)
}
```

### Rename

This option is disabled by default. 

It is often the case that functions passed to dig as a provider are private functions not visible in the main package. In this case, you can enable this parameter. Then anti-dig will not only generate code but also change the names of used functions, making them public and available in the main

```go
package main

import dig "github.com/3timeslazy/anti-dig"

func main() {
	dig.Anti.Rename(true)
}
```

