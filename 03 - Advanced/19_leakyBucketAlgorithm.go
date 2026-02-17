package main

import (
	"math"
	"sync"
	"time"
)

type LeakyBucket struct {
	capacity int
	leakRate time.Duration
	tokens   int
	lastLeak time.Time
	mu       sync.Mutex
}

func NewLeakyBucket(capacity int, leakRate time.Duration) *LeakyBucket {
	return &LeakyBucket{
		capacity: capacity,
		leakRate: leakRate,
		tokens:   capacity,
		lastLeak: time.Now(),
	}
}

func (b *LeakyBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsedTime := now.Sub(b.lastLeak)
	tokensToAdd := int(elapsedTime / b.leakRate)
	if tokensToAdd > 0 {
		b.tokens = int(math.Min(float64(b.capacity), float64(b.tokens+tokensToAdd)))
		b.lastLeak = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

func main() {
	leakyBucketInstance := NewLeakyBucket(5, 500*time.Millisecond)

	numRequests := 10
	for range numRequests {
		if leakyBucketInstance.Allow() {
			println("Request allowed")
		} else {
			println("Request denied")
		}
		time.Sleep(150 * time.Millisecond)
	}
}
