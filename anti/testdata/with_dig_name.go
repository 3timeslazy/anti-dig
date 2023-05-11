package main

import (
	"github.com/3timeslazy/anti-dig/example/consumer"
	"github.com/3timeslazy/anti-dig/example/consumer/queue"
)

func Provide() consumer.Consumer {
	var2_queue1 := queue.New1()
	var3_queue2 := queue.New2()
	var4 := consumer.ConsumerParams{
		Queue1:	var2_queue1,
		Queue2:	var3_queue2,
	}
	var1 := consumer.New(var4)
	return var1
}
