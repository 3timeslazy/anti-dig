package server

import dig "github.com/3timeslazy/anti-dig"

type RunParams struct {
	dig.In
	Config Config
}

func Run(RunParams) {}

type Config struct{}

func NewConfig() Config {
	return Config{}
}
