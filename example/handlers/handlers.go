package handlers

import (
	dig "github.com/3timeslazy/anti-dig"
)

type HandlersParam struct {
	dig.Out
	Handler Handler `group:"handlers"`
}

type HandlersFlattenParam struct {
	dig.Out
	Handlers []Handler `group:"handlers,flatten"`
}

type Handler interface {
	Handle() error
}
