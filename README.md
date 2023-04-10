# :carpentry_saw: anti-dig

A 100% compatible drop-in replacement of `go.uber.org/dig`

**anti-dig** does the same thing as dig, but instead of calling providers and passing them on to other providers in runtime, it generates a file with the code calling your providers in the correct order

# Usage

Replace `go.uber.org/dig` with `github.com/3timeslazy/anti-dig`
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

and run your code. It'll generate a new file with all dependencies in one file

# Example
You can find the code [here](https://github.com/3timeslazy/anti-dig/blob/main/example/main.go)

The code below 
```go
package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/cron"
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/handlers/flatten"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
	"github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
	"github.com/3timeslazy/anti-dig/example/observability"
	"github.com/3timeslazy/anti-dig/example/server"
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
	err = container.Provide(handlerv0.NewHandlerV0)
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
	err = container.Provide(server.NewServer)
	if err != nil {
		panic(err)
	}
	err = container.Invoke(Run)
	if err != nil {
		panic(err)
	}
}

func Run(cron cron.Cron, server *server.Server) {
	go cron.Start()
	go server.Serve()
}
```

will generate the following
```go
package main

import (
        db "github.com/3timeslazy/anti-dig/example/db"
        config "github.com/3timeslazy/anti-dig/example/config"
        cron "github.com/3timeslazy/anti-dig/example/cron"
        handlers "github.com/3timeslazy/anti-dig/example/handlers"
        server "github.com/3timeslazy/anti-dig/example/http/server"
        observability "github.com/3timeslazy/anti-dig/example/observability"
        flatten "github.com/3timeslazy/anti-dig/example/handlers/flatten"
        handlerv0 "github.com/3timeslazy/anti-dig/example/handlers/handlerv0"
        handlerv1 "github.com/3timeslazy/anti-dig/example/handlers/handlerv1"
)

func main() {
	var2, err := db.NewDB()
	if err != nil {
		panic(err)
	}
	var3 := config.NewConfig()
	var1 := cron.NewCron(var2, var3)
	var6_0 := observability.NewObservability(var3)

	var8_0 := flatten.NewListOfHandlers(var6_0.Metrics)

	var10_0, err := handlerv0.NewHandlerV0(var2)
	if err != nil {
		panic(err)
	}

	var10_1 := handlerv1.NewHandlerV1()

	var11 := []handlers.Handler{
		var10_0.Handler,
		var10_1.Handler,
	}
	var11 = append(var11, var8_0.Handlers...)
	var12 := server.ServerParams{
		Handlers: var11,
	}
	var4 := server.NewServer(var12)
	FIXME(var1, var4)
}
```

After that, simply replace the `main()` function in the original file with the generated function, replace `FIXME()` with your function, remove `go.uber.org/dig` from dependencies and you're all set! You can now run the generated code, where you can <ins>clearly see all the dependencies between types<ins>.

But before that, I'd suggest renaming the variables and maybe grouping the code better. ChatGPT it's pretty good at this, but I wouldn't recommend using it with proprietary code
