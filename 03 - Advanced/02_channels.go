package main

import "fmt"

func main() {
	
	chanA := make(chan string);
	go func () {
		for _, e := range "abcde" {
			chanA <- "Alphabet: " + string(e) 
		}
	}()
	for range 5 {
		receiver := <- chanA
		fmt.Println(receiver);
	}
}