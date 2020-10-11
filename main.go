package main

import (
	"fmt"
	"messageQueue/bus"
	"messageQueue/sub"
	"os"
	"time"
)

func main() {
	// Server starts working
	bus.Start()

	var exchangeName = "superman"

	bus.AddExchange(exchangeName)

	var subers  = make([]sub.Subscriber, 0, 10)
	for i := 0; i < cap(subers); i++ {
		var _sub sub.Subscriber
		subers = append(subers, _sub)
	}

	bus.AddSubsToExchange(exchangeName, subers...)

	// Publish message
	for {
		if err := bus.Publish([]byte("[I'm publishing]: hello world"), exchangeName); err != nil {
			fmt.Println(err)
			os.Exit(127)
		}
		time.Sleep(time.Millisecond * 500)
	}

}
