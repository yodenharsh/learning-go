package main

import "fmt"

func main() {

	var hello string = "hello there"
	anotherString := "hi there"

	// Formatting
	s := fmt.Sprint(hello, anotherString, 124, 52)
	fmt.Println(s)

	s = fmt.Sprintln(hello, anotherString, 124, 52)
	fmt.Println(s)

	// Input scanning
	var name string
	var age int

	fmt.Print("Enter name and age: ")
	fmt.Scan(&name, &age)

	fmt.Printf("Name: %s, Age: %d\n", name, age)

	err := checkAge(15)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func checkAge(age int) error {
	if age < 18 {
		return fmt.Errorf("%d is not legal", age)
	}
	return nil
}
