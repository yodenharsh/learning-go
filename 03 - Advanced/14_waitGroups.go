package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	// wg.Add(1) Could also just increment it here instead but this is generally worse practice as
	// it can lead to race conditions if the worker is launched before the main goroutine has a chance to increment the counter
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second) // Simulate work
	fmt.Printf("Worker %d done\n", id)
}

func workerWithChannel(id int, result chan<- int) {
	fmt.Printf("Worker %d starting\n", id)
	time.Sleep(time.Second)

	result <- id * 10 // Send result back to main goroutine
	fmt.Printf("Worker %d done\n", id)
}

type ConstructionWorker struct {
	Id   int
	Task string
}

func (w *ConstructionWorker) PerformTask() {
	fmt.Printf("Worker %d performing task: %s\n", w.Id, w.Task)
	time.Sleep(1 * time.Second)
	fmt.Printf("Worker ID %d completed task: %s\n", w.Id, w.Task)
}

func main() {
	// Scenario 1
	var wg sync.WaitGroup
	numWorkers := 3

	// Note: instead of manually adding and calling Done() later, we can simply use wg.Go() too
	// which will handle adding and defering the Done() call for us, but this is just to illustrate the basic usage of WaitGroups
	wg.Add(numWorkers)
	// Launching workers
	for i := range numWorkers {
		go worker(i, &wg)
	}

	wg.Wait()
	fmt.Println("All workers completed")

	// Scenario 2 (with channels)
	resultChan := make(chan int, numWorkers) // Buffered channel to avoid blocking
	for i := range numWorkers {
		wg.Go(func() { workerWithChannel(i, resultChan) })
	}

	go func() {
		wg.Wait() // Wait for all workers to finish before closing the channel
		close(resultChan)
	}()

	for result := range resultChan {
		fmt.Printf("Received result: %d\n", result)
	}

	// Scenario 3 (with struct and method)
	tasks := []string{"Build foundation", "Erect walls", "Install roof"}
	for i, task := range tasks {
		worker := ConstructionWorker{Id: i, Task: task}
		wg.Go(worker.PerformTask)
	}

	wg.Wait()
	fmt.Println("Construction completed")
}
