package main

import (
	"fmt"
	"sync"
)

var once sync.Once

func initialize() {
	fmt.Println("This will not be repeated even if called multiple times.")
}

func main() {

	var wg sync.WaitGroup
	for i := range 5 {
		fmt.Println("Called initialize function")
		fmt.Println("Running goroutine #", i)
		wg.Go(func() { once.Do(initialize) })
	}

	wg.Wait()
}
