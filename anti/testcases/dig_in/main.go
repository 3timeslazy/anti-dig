package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_in/server"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(server.NewConfig)

	_ = container.Invoke(server.Run)
}
