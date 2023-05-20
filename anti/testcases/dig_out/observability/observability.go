package observability

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_out/metrics"
)

func New() Observability {
	return Observability{}
}

type Observability struct {
	dig.Out
	Metrics metrics.Metrics
	Profiler
}

type Profiler struct{}
