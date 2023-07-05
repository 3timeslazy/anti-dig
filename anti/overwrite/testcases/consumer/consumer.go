package consumer

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/overwrite/testcases/registry"
)

type Consumer struct{}

func init() {
	registry.Provide(newConsumer)
}

type consumerParam struct {
	dig.In
}

func newConsumer(consumerParam) Consumer {
	return Consumer{}
}
