package handlers

import (
	dig "github.com/3timeslazy/anti-dig"
)

type HandlersGRPC struct {
	dig.Out
	Handler Handler `group:"grpc_handlers"`
}

type HandlersHTTP struct {
	dig.Out
	Handler Handler `group:"http_handlers"`
}

type HandlersFlattenParam struct {
	dig.Out
	Handlers []Handler `group:"http_handlers,flatten"`
}

type Handler interface {
	Handle() error
}
