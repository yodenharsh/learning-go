package main

import "fmt"

// Note: "send" and "receive" from our perspective and not channel's

func main() {
	chanA := make(chan int)

	go func(ch chan<- int) { // This is a send only channeel
		for i := range 5 {
			ch <- i
		}
		close(ch)
	}(chanA)

	for value := range chanA {
		fmt.Println("Received: ", value)
	}

	chanB := make(chan int)
	receiveData(chanB)
}

func receiveData(ch <-chan int) {
	for value := range ch {
		fmt.Println("Received: ", value)
	}
}
