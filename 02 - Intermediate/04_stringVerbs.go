package main

import "fmt"

func main() {
	aString := "Hello there"

	fmt.Printf("%v\n", aString)
	fmt.Printf("As it appears in Go: %#v\n", aString)
	fmt.Printf("%T\n", aString)
}
