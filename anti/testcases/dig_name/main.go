package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_name/handlers"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(handlers.New, dig.Name("handler_v1"))
	_ = container.Provide(handlers.New, dig.Name("handler_v2"))
	_ = container.Provide(handlers.NewV3, dig.Name("handler_v3"))
	_ = container.Provide(handlers.NewServer)

	_ = container.Invoke(func(handlers.Server) {})
}
