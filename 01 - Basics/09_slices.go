package main

import (
	"fmt"
	"slices"
)

func main() {

	numbers := [5]int{1, 5, 1, 51, 12}

	slice := make([]int, 5)
	var _ = slice

	numbersSlice := numbers[1:4]
	fmt.Println("numbers[1:4] = ", numbersSlice)

	numbersSlice = append(numbersSlice, 19)
	fmt.Println("After appending 19 = ", numbersSlice)
	fmt.Println("Original array = ", numbers)

	numbersSlice = append(numbersSlice, 185, 45, 1543)
	fmt.Println("appending to overflow = ", numbersSlice)
	fmt.Println("Original array = ", numbers)

	numbersSliceCopy := make([]int, len(numbersSlice))
	copy(numbersSliceCopy, numbersSlice)

	if slices.Equal(numbersSlice, numbersSliceCopy) {
		fmt.Println("Slices are equal and are equated by values rather than memory reference")
	} else {
		fmt.Println("Slices are unequal")
	}
}
