package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	// for tick := range ticker.C {
	// 	fmt.Println("Ticking at: ", tick)
	// }

	ticker2 := time.NewTicker(2 * time.Second)
	defer ticker2.Stop()

	stop := time.After(7 * time.Second)

	for {
		select {
		case <-ticker2.C:
			periodicTask()
		case <-stop:
			fmt.Println("Stopping ticker after 7 seconds")
			return
		}
	}
}

func periodicTask() {
	fmt.Println("Performing periodic task at: ", time.Now())
}
