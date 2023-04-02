package cron

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/db"
)

type Cron struct{}

func NewCron(bny db.DB, cny config.Config) Cron {
	return Cron{}
}
