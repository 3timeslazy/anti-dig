package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_out/metrics"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_out/observability"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(observability.New)
	_ = container.Invoke(func(metrics.Metrics) {})
}
