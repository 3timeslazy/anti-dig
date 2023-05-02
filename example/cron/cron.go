package cron

import (
	"github.com/3timeslazy/anti-dig/example/config"
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/db"
)

type Cron interface {
	Start()
}

func NewCron(db db.DB, consumer consumer.Consumer, config config.Config) Cron {
	return nil
}
