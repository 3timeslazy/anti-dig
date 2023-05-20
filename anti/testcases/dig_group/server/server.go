package server

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/dig_group/handlers"
)

type ServerParams struct {
	dig.In
	HandlersV1 []handlers.Handler `group:"handlers_v1"`
	HandlersV2 []handlers.Handler `group:"handlers_v2"`
}

type Server struct{}

func New(params ServerParams) Server {
	return Server{}
}
