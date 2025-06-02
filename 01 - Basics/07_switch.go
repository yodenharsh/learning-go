package main

import "fmt"

func main() {

	const fruit = "Apple"

	switch fruit {
	case "Apple":
		fmt.Println("This is an apple")
	case "Banana":
		fmt.Println("Banana")
	}

	switch {
	case 10 > 15:
		fmt.Println("This is not going to execute")
		fallthrough
	case 15 == 15:
		fmt.Println("Hello there")
		fmt.Println("Or not Ig")
	default:
		fmt.Println("This should not execute either")
	}

	checkType(12)
}

func checkType(x interface{}) {
	switch x.(type) {
	case int:
		fmt.Println("Integer")
	case float64:
		fmt.Println("It is a float")
	case string:
		fmt.Println("String it is")
	default:
		fmt.Println("Unknown type")
	}
}
