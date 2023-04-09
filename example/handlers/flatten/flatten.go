package flatten

import (
	"github.com/3timeslazy/anti-dig/example/handlers"
	"github.com/3timeslazy/anti-dig/example/observability"
)

func NewListOfHandlers(_ observability.Metrics) handlers.HandlersFlattenParam {
	return handlers.HandlersFlattenParam{
		Handlers: []handlers.Handler{
			HandlerA{},
			HandlerB{},
		},
	}
}

type HandlerA struct{}

func (h HandlerA) Handle() error { return nil }

type HandlerB struct{}

func (h HandlerB) Handle() error { return nil }
