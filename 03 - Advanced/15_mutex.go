package main

import (
	"fmt"
	"sync"
)

type counter struct {
	mu    sync.Mutex
	value int
}

func (c *counter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *counter) getValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	var wg sync.WaitGroup
	c := &counter{}

	numGoroutines := 10
	for range numGoroutines {
		wg.Go(func() {
			for range 1000 {
				c.increment()
				// c.value++ // This is not thread-safe and
			}
		})
	}
	wg.Wait()
	fmt.Println("Final counter value:", c.getValue())
}
