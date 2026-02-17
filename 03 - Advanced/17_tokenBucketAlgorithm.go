package main

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens     chan struct{}
	refillTime time.Duration
}

func NewRateLimiter(rateLimit int, refillTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		tokens:     make(chan struct{}, rateLimit),
		refillTime: refillTime,
	}

	for range rateLimit {
		rl.tokens <- struct{}{}
	}

	go rl.StartRefill()

	return rl
}

func (rl *RateLimiter) Allow() bool {
	select {
	case <-rl.tokens:
		return true
	default:
		return false
	}
}

func (rl *RateLimiter) StartRefill() {
	ticker := time.NewTicker(rl.refillTime)
	defer ticker.Stop()

	for range ticker.C {
		rl.tokens <- struct{}{}
	}
}

func main() {
	rateLimiter := NewRateLimiter(5, time.Second)

	numRequests := 10
	for range numRequests {
		if rateLimiter.Allow() {
			fmt.Println("Request allowed")
		} else {
			fmt.Println("Request denied")
		}
		time.Sleep(300 * time.Millisecond)
	}
}
