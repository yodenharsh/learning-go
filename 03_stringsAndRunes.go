package main

import "fmt"

func main() {
	simpleString := "Hello\nNewline"
	rawString := `Hello\nNewline`

	fmt.Println("Simple string: ", simpleString)
	fmt.Println("Raw string: ", rawString)

	fmt.Println("Simple string length: ", len(simpleString))

	fmt.Println("No need to do any localsCompare for string comparison. Using normal relational ops is fine. ", "apple" > "banana")

	for _, char := range rawString {
		fmt.Printf("C: %c, ", char)
	}
	fmt.Println()

	// Runes have to be declared with a single quote (')
	var aRune = 'a'
	fmt.Println("Runes need to be declared with a single quote '", aRune)

	// Converting a rune to a string
	var _ = string(aRune)

}
