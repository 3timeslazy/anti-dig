package server

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/example/handlers"
)

type Server struct{}

type ServerParams struct {
	dig.In
	Handlers []handlers.Handler `group:"handlers"`
}

func NewServer(_ ServerParams) *Server {
	return nil
}
