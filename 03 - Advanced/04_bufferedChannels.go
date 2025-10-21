package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting scenario 1");
	scenario1();

	fmt.Println("Starting scenario 2");
	scenario2();
}

func scenario1 () {
	ch := make(chan int, 2);
	
	ch <- 1
	ch <- 2
	// ch <- 3; This breaks because now the channel is looking for receivers before accepting the new value
	
	funcA := func () {
		time.Sleep(2 * time.Second);
		fmt.Println("Received: ", <- ch);
		fmt.Println("Blocking ends")
	}
	go funcA()
	
	ch <- 3
	fmt.Println("Recieved out: ", <- ch);
	fmt.Println("Recieved out: ", <- ch);

	fmt.Println("buffered channels");
}

func scenario2() {
	ch := make (chan int, 2);

	go func() {
		time.Sleep(2 * time.Second);
		ch <- 1;
		ch <- 2;
	}()

	fmt.Println("Value: ", <- ch);
	fmt.Println("Value: ", <- ch);
	fmt.Println("Ending")
}