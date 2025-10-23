package main

import (
	"fmt"
	"time"
)

func main() {

	// Simple closing channel example
	ch := make(chan int)

	go func() {
		for i := range 5 {
			ch <- i
		}
		close(ch)
	}()

	for val := range ch {
		fmt.Println(val)
	}

	// Receiving from a closed channel
	ch2 := make(chan int)
	close(ch2)

	val, ok := <-ch2
	if !ok {
		fmt.Println("Channel is closed")
	} else {
		fmt.Println("Impossible to reach here: ", val)
	}

	// Range over closed channel
	ch3 := make(chan int)
	go func() {
		for i := range 5 {
			ch3 <- i
		}
		close(ch3)
	}()

	for val := range ch3 {
		fmt.Println("Received: ", val)
	}

	// Runtime panic when closnig channels twice
	ch4 := make(chan int)
	go func() {
		close(ch4)
		// close(ch4) This causes a panic
	}()

	time.Sleep(time.Second)

	// pipeline pattern
	ch5 := make(chan int)
	ch6 := make(chan int)
	go producer(ch5)
	go filter(ch5, ch6)

	for val := range ch6 {
		fmt.Println("Value from ch6: ", val)
	}

}

func producer(ch chan<- int) {
	for i := range 5 {
		ch <- i
	}
	close(ch)
}

func filter(in <-chan int, out chan<- int) {
	for val := range in {
		if val%2 == 0 {
			out <- val
		}
	}

	close(out)
}
