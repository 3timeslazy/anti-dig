package main

import (
	"github.com/3timeslazy/anti-dig/example/observability"
	"github.com/3timeslazy/anti-dig/example/config"
)

func Provide() observability.Metrics {
	var2 := config.NewConfig()
	var3_0 := observability.NewObservability(var2)

	return var3_0.Metrics
}
