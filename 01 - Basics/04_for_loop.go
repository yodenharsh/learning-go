package main

import "fmt"

func main() {
	for i := 1; i <= 5; i++ {

		fmt.Printf("i's value = %d\n", i)
	}

	numbers := []int{1, 2, 3, 4, 5, 6}
	for index, number := range numbers {
		fmt.Printf("Number = %v at index %v\n", number, index)
	}
}
