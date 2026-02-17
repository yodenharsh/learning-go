package main

import "time"

type StatefulWorker struct {
	count int
	ch    chan int
}

func (w *StatefulWorker) Start() {
	go func() {
		for value := range w.ch {
			w.count += value
			println("Current count:", w.count)

		}
	}()
}

func (w *StatefulWorker) Send(value int) {
	w.ch <- value
}

func main() {
	stWorker := &StatefulWorker{
		count: 0,
		ch:    make(chan int),
	}
	stWorker.Start()
	for i := range 5 {
		stWorker.Send(i)
		time.Sleep(500 * time.Millisecond)
	}
}
