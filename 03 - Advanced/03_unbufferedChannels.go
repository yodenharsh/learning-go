package main

import "fmt"

func main () {
	ch := make(chan int)

	go func(){
	ch <- 1;
	}();

	receiver := <- ch
	fmt.Println(receiver);
}