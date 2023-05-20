package handlers

import dig "github.com/3timeslazy/anti-dig"

type Handler struct{}

type HandlersIn struct {
	dig.In
	H1 Handler `name:"handler_v1"`
	// H2 Handler `name:"handler_v2"`
	H3 Handler `name:"handler_v3"`
}

func New() Handler {
	return Handler{}
}

func NewV3() Handler {
	return Handler{}
}

type Server struct{}

func NewServer(HandlersIn) Server {
	return Server{}
}
