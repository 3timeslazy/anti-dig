package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_group/handlers"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_group/server"
)

func main() {
	// The case checks how anti-dit works with:
	// - the `group` tag
	// - the `flatten` tag
	// - two different groups at the same time
	// - dig.Group(...)

	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(server.New)
	_ = container.Provide(handlers.NewListV1)
	_ = container.Provide(handlers.NewListV2)
	_ = container.Provide(handlers.NewV1)
	_ = container.Provide(handlers.NewV2, dig.Group("handlers_v2"))

	_ = container.Invoke(func(server.Server) {})
}
