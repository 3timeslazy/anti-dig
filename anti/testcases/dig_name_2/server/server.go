package server

import dig "github.com/3timeslazy/anti-dig"

type Handler struct{}

type HandlerV1 struct {
	dig.Out
	Handler *Handler `name:"v1"`
}

func NewV1() HandlerV1 {
	return HandlerV1{}
}

type HandlerV2 struct {
	dig.Out
	Handler *Handler `name:"v2"`
}

func NewV2() HandlerV2 {
	return HandlerV2{}
}

type Handlers struct {
	dig.In
	V1 *Handler `name:"v1"`
	V2 *Handler `name:"v2"`
}

type Server struct{}

func NewServer(Handlers) Server {
	return Server{}
}
