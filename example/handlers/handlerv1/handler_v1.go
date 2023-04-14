package handlerv1

import "github.com/3timeslazy/anti-dig/example/handlers"

func NewHandlerV1() handlers.HandlersGRPC {
	return handlers.HandlersGRPC{
		Handler: HandlerV1{},
	}
}

type HandlerV1 struct{}

func (h HandlerV1) Handle() error { return nil }
