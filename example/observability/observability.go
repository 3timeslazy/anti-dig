package observability

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/example/config"
)

func NewObservability(_ config.Config) Observability {
	return Observability{}
}

type Observability struct {
	dig.Out
	Metrics
	Profiler Profiler
}

type Metrics struct{}

func NewMetrics() Metrics {
	return Metrics{}
}

type Profiler struct{}

func NewProfiler() Profiler {
	return Profiler{}
}
