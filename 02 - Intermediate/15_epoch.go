package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	unixTime := now.Unix()

	fmt.Println("Now in Unix: ", unixTime)

	t := time.Unix(unixTime, 0)
	fmt.Println(t)
}
