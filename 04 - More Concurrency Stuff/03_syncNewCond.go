package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const bufferSize = 5

type Buffer struct {
	items []int
	mu    sync.Mutex
	cond  *sync.Cond
}

func newBuffer(size int) *Buffer {
	buffer := &Buffer{items: make([]int, 0)}
	buffer.cond = sync.NewCond(&buffer.mu)

	return buffer
}

func (b *Buffer) produce(item int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.items) == bufferSize {
		b.cond.Wait()
	}
	b.items = append(b.items, item)
	fmt.Println("Produced item: ", item)

	b.cond.Signal()
}

func (b *Buffer) consume() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	for len(b.items) == 0 {
		b.cond.Wait()
	}

	item := b.items[0]
	b.items = b.items[1:]
	fmt.Println("Consumed item: ", item)

	b.cond.Signal()
	return item
}

func producerHelper(buffer *Buffer) {
	for i := range 10 {
		buffer.produce(i * rand.Intn(20))
		time.Sleep(500 * time.Millisecond)
	}
}

func consumerHelper(buffer *Buffer) {
	for range 10 {
		buffer.consume()
		time.Sleep(700 * time.Millisecond)
	}
}

func main() {
	buffer := newBuffer(bufferSize)
	var wg sync.WaitGroup

	wg.Go(func() { producerHelper(buffer) })
	wg.Go(func() { consumerHelper(buffer) })

	wg.Wait()
}
