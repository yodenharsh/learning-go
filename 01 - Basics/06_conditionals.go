package main

import "fmt"

func main() {
	var age int

	fmt.Scanln(&age)
	if age >= 18 {
		fmt.Printf("You may participate in this thing")
	} else {
		fmt.Printf("Lmao too young")
	}

	if true {

	} else {
		fmt.Printf("This is non-reachable code")
	}
}
