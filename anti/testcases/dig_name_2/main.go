package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_name_2/server"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(server.NewServer)
	_ = container.Provide(server.NewV1)
	_ = container.Provide(server.NewV2)

	_ = container.Invoke(func(server.Server) {})
}
