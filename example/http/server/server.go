package server

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/handlers"
)

type Server struct{}

type ServerParams struct {
	dig.In
	Config   config.Config
	Handlers []handlers.Handler `group:"http_handlers"`
}

func NewServer(_ ServerParams) *Server {
	return nil
}

func (srv *Server) Serve() {}
