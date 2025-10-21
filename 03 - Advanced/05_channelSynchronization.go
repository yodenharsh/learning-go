package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	done := make(chan int)

	go func() {
		fmt.Println("Working...")
		time.Sleep(2 * time.Second)
		done <- 0
	}()

	<-done
	fmt.Println("Finished.")

	ch := make(chan int)

	go func() {
		ch <- 9
		fmt.Println("Sent value")
	}()

	value := <-ch
	time.Sleep(1 * time.Second)
	fmt.Println(value)

	numGoRoutines := 3
	done = make(chan int, 3)

	for i := range numGoRoutines {
		go func(id int) {
			fmt.Printf("Goroutine %d working\n", id)
			done <- id
		}(i)
	}

	for range numGoRoutines {
		<-done // Wait for each goroutine to finish
	}

	fmt.Println("All goroutines are finished")

	data := make(chan string)

	go func() {
		for i := range numGoRoutines {
			data <- "Hello " + strconv.Itoa(i)
			time.Sleep(100 * time.Millisecond)
		}
		close(data)
	}()

	for value := range data {
		fmt.Println("Received value: ", value, ":", time.Now())
	}
}
