package main

import "fmt"

func main() {

	sum := 0
	for sum < 15 {
		sum += 3
		fmt.Println("Sum = ", sum)
	}
}
