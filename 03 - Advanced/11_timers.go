package main

import (
	"fmt"
	"time"
)

func main() {

	timer := time.NewTimer(2 * time.Second)

	isStoppedSuccess := timer.Stop()
	if isStoppedSuccess {
		fmt.Println("Timer stopped successfully")
	} else {
		<-timer.C // Blocking until the timer expires
	}
	// timer.Reset(3 * time.Second)
	fmt.Println("Timer expired")

	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	go func() {
		longRunningTask()
		done <- true
	}()

	select {
	case <-timeout:
		fmt.Println("Timeout occurred")
	case <-done:
		fmt.Println("Long running task completed")
	}

	timer = time.NewTimer(2 * time.Second)
	// THIS WILL NOT BLOCK and program will end
	go func() {
		<-timer.C
		fmt.Println("Delayed operation executed")
	}()

	fmt.Println("We are waiting")
}

func longRunningTask() {
	fmt.Println("Starting long running task...")
	for i := range 20 {
		fmt.Println(i)
		time.Sleep(time.Second)
	}
}
