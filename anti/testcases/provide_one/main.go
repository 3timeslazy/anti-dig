package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/provide_one/config"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(config.New)
	_ = container.Invoke(func(config.Config) {})
}
