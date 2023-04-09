package cron

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/db"
)

type Cron struct{}

func NewCron(db db.DB, config config.Config) Cron {
	return Cron{}
}

func (cron *Cron) Start() {}
