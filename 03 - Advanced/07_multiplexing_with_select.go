package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		time.Sleep(time.Second)
		ch1 <- 1

	}()

	go func() {
		time.Sleep(time.Second)
		ch2 <- 5
	}()

	time.Sleep(time.Millisecond * 1500)

	for range 2 {
		select {
		case msg := <-ch1:
			fmt.Println("Received from channel 1: ", msg)

		case msg := <-ch2:
			fmt.Println("Received from channel 2: ", msg)

			// default:
			// 	fmt.Println("No channels ready")
		}
	}

	ch := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		ch <- 1
	}()

	select {
	case msg := <-ch:
		fmt.Println("Received: ", msg)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

	demoOk()
}

func demoOk() {
	ch := make(chan int)

	go func() {
		ch <- 1
		close(ch)
	}()

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed")
				return
			} else {
				fmt.Println(msg)
			}
		}
	}
}
