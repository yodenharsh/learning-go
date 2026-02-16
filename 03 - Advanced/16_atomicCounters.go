package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type AtomicCounter struct {
	count int64
}

func (ac *AtomicCounter) Increment() {
	atomic.AddInt64(&ac.count, 1)
}

func (ac *AtomicCounter) getValue() int64 {
	return atomic.LoadInt64(&ac.count)
}

func main() {
	var wg sync.WaitGroup
	numRoutines := 100
	counter := &AtomicCounter{}

	for range numRoutines {
		wg.Go(func() {
			for range 1000 {
				counter.Increment()
			}
		})
	}
	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.getValue())
}
