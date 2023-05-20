package main

import (
	dig "github.com/3timeslazy/anti-dig"
	"github.com/3timeslazy/anti-dig/anti/testcases/provide_many/config"
	"github.com/3timeslazy/anti-dig/anti/testcases/provide_many/db"
)

func main() {
	dig.Anti.Optimise(false)

	container := dig.New()
	_ = container.Provide(config.New)
	_ = container.Provide(db.New)
	_ = container.Invoke(func(db.DB) {})
}
