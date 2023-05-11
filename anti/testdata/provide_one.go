package main

import (
	"github.com/3timeslazy/anti-dig/example/config"
)

func Provide() config.Config {
	var1 := config.NewConfig()
	return var1
}
