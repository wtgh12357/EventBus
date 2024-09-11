package main

import (
	"fmt"
	"time"
)

func main() {

	bus := NewEventBus()
	subscribe := bus.Subscribe("test")
	go func() {
		fmt.Println("写入")
		subscribe <- Event{Payload: 1}
	}()
	time.Sleep(2 * time.Second)
	e := <-subscribe
	fmt.Println(e.Payload)
}
