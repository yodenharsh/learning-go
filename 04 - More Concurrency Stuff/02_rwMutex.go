package main

import (
	"math/rand"
	"sync"
)

type Counter struct {
	counter int
	mu      sync.RWMutex
}

func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counter++
}

func (c *Counter) Get() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.counter
}

func main() {

	var wg sync.WaitGroup

	counter := &Counter{}
	for range 20000 {
		randomNum := rand.Intn(100)
		if randomNum > 90 {
			wg.Go(counter.Increment)
		}

		wg.Go(func() { counter.Get() })
	}

	wg.Wait()
}
