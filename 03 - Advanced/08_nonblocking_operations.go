package main

import (
	"fmt"
	"time"
)

func main() {
	// Non blocking receive operation
	ch := make(chan int)
	select {
	case val := <-ch:
		fmt.Println("Got ch: ", val)
	default:
		fmt.Println("No message available")
	}

	// Non blocking send operation
	select {
	case ch <- 1:
		fmt.Println("Sent message")
	default:
		fmt.Println("Channel is not ready to receive")
	}

	// Non blocking operation in realtime systems
	nonBlockingRealTimeDemo()

}

func nonBlockingRealTimeDemo() {
	data := make(chan int)
	quit := make(chan bool)

	go func() {
		for {
			select {
			case d := <-data:
				fmt.Println("Data received: ", d)

			case <-quit:
				fmt.Println("Stopping")
				return

			default:
				fmt.Println("Waiting for data")
				time.Sleep(500 * time.Millisecond)
			}
		}
	}()

	for i := range 5 {
		data <- i
		time.Sleep(time.Second)
	}

	quit <- true
}
