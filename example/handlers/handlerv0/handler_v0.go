package handlerv0

import (
	"github.com/3timeslazy/anti-dig/example/db"
	"github.com/3timeslazy/anti-dig/example/handlers"
)

func NewHandlerV0(_ db.DB) (handlers.Handler, error) {
	return HandlerV0{}, nil
}

type HandlerV0 struct{}

func (h HandlerV0) Handle() error { return nil }
