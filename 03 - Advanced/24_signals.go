package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pid := os.Getpid()
	fmt.Printf("Process ID: %d\n", pid)

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs

		switch sig {
		case syscall.SIGINT:
			fmt.Println("Received SIGINT signal")
		case syscall.SIGTERM:
			fmt.Println("Received SIGTERM signal")
		default:
			fmt.Println("Received unknown signal: ", sig)
		}
		os.Exit(0)
	}()

	// Simulate long-running process
	fmt.Println("Running... Press Ctrl+C to exit.")
	time.Sleep(200 * time.Second)

}
