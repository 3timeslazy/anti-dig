package consumer

import dig "github.com/3timeslazy/anti-dig"

type ConsumerParams struct {
	dig.In
	Queue1 Queue `name:"queue_1"`
	Queue2 Queue `name:"queue_2"`
}

type Queue struct{}

type Consumer struct{}

func New(params ConsumerParams) Consumer {
	return Consumer{}
}
