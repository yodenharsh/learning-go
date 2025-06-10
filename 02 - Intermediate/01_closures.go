package main

import "fmt"

func main() {
	sequence := adder()

	sequence(15)
	sequence(30)

	// This will use a difference memory than adder
	sequence2 := adder()
	sequence2(12)
	sequence2(21)
}

func adder() func(someValue int) int {
	i := 0

	fmt.Println("Previous value of i: ", i)
	return func(someValue int) int {
		i = i + someValue
		fmt.Println("New i value: ", i)
		return i
	}
}
