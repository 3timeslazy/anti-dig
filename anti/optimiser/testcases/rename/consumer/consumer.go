package consumer

import (
	"github.com/3timeslazy/anti-dig/anti/optimiser/testcases/rename/registry"
)

type Consumer struct{}

func init() {
	registry.Provide(newConsumer)
}

func newConsumer() Consumer {
	return Consumer{}
}
