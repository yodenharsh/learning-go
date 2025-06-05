package main

import (
	"fmt"
)

func main() {

	var numbers [5]int
	fmt.Println("Printing numbers uninitialized: ", numbers)

	numbers[3] = 12

	fruits := [5]string{"omg", "how", "is", "this", "readable"}

	for _, fruit := range fruits {
		fmt.Println(fruit)
	}

}
