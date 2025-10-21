package main

import (
	"fmt"
	"time"
)

func main () {
var error error;

	go hello();

	go printNumbers();
	go printLetters();

	go func() {
		error = doWork();
	}();

	
	time.Sleep(2 * time.Second);
	if error != nil {
		fmt.Println("Error: ", error);
	} else {
		fmt.Println("Work completed successfully");
	}
}

func hello () {
	time.Sleep(1 * time.Second);
	fmt.Println("Hello from Goroutine");
}

func printNumbers() {
	for range 5 {
		fmt.Println(time.Now());
		time.Sleep(100 * time.Millisecond);
	}
}

func printLetters () {
	for _, letter := range "abcdefghijklmnopqrstuvxyz" {
		fmt.Println(string(letter) + " with time: " + string(time.Now().GoString()));
		time.Sleep(200 * time.Millisecond);
	}
}

func doWork() error {
	time.Sleep(1*  time.Second);
	return fmt.Errorf("An error occurred somehow.")
}