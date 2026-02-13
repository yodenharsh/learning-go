package main

import (
	"fmt"
	"time"
)

func worker(id int, tasks <-chan int, results chan<- int) {
	for task := range tasks {
		fmt.Printf("Worker %d processing task %d\n", id, task)
		// Simulating work
		time.Sleep(1 * time.Second)
		results <- task * 2 // Just an example of processing
	}
}

type ticketRequest struct {
	personId   int
	numTickets int
	cost       int
}

func ticketProcessor(requests <-chan ticketRequest, results chan<- int) {
	for request := range requests {
		fmt.Printf("Processing %d ticket(s) for personId %d with total cost %d\n", request.numTickets, request.personId, request.cost)
		// Simulating processing time
		time.Sleep(500 * time.Millisecond)
		results <- request.personId
	}
}

func main() {
	// Scenario 1
	numWorkers := 3
	numJobs := 10

	tasks := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Creating workers
	for i := range numWorkers {
		go worker(i, tasks, results)
	}

	// Sending values to the tasks channel
	for i := range numJobs {
		tasks <- i
	}
	close(tasks)

	for range numJobs {
		result := <-results
		fmt.Printf("Result: %d\n", result)
	}

	// Secnario 2
	numRequests := 5
	price := 10
	ticketRequests := make(chan ticketRequest, numRequests)
	ticketResults := make(chan int, numRequests)

	// Creating ticket processor
	for range 2 {
		go ticketProcessor(ticketRequests, ticketResults)
	}

	for i := range numRequests {
		ticketRequests <- ticketRequest{
			personId:   i,
			numTickets: i*2 + 1,
			cost:       (i + 1) * price,
		}
	}

	close(ticketRequests)

	for range numRequests {
		result := <-ticketResults
		fmt.Printf("Ticket processed for personId: %d\n", result)
	}
}
