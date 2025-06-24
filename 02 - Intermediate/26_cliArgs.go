package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Command: ", os.Args[0])
	for i, arg := range os.Args {
		fmt.Println("Argument: ", i, ": ", arg)
	}

	// Define flags
	var name string

	var age int
	var isMale bool

	flag.StringVar(&name, "name", "John", "Name of user")
	flag.IntVar(&age, "age", 18, "Age of user")
	flag.BoolVar(&isMale, "male", true, "Gender of user")

	flag.Parse()
	fmt.Println("Name: ", name)
	fmt.Println("Age: ", age)
	fmt.Println("Is male: ", isMale)
}
