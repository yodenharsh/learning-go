package main

import (
	"sync"
	"time"
)

type RateLimiter struct {
	mu     sync.Mutex
	count  int
	limit  int
	window time.Duration
	reset  time.Time
}

func newRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if now.After(rl.reset) {
		rl.count = 0
		rl.reset = now.Add(rl.window)
	}
	if rl.count < rl.limit {
		rl.count++
		return true
	}

	return false
}

func main() {
	rateLimiter := newRateLimiter(5, time.Second)

	for range 10 {
		if rateLimiter.Allow() {
			println("Request allowed")
		} else {
			println("Request denied")
		}
		time.Sleep(150 * time.Millisecond)
	}
}
