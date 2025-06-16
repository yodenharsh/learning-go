package main

import "fmt"

type PhoneNumber struct {
	countryCode int
	cell        string
}

type Person struct {
	firstName string
	lastName  string
	age       int
	PhoneNumber
}

func (p Person) fullName() string {
	return fmt.Sprintf("%s %s", p.firstName, p.lastName)
}

func (p *Person) add1ToAge() {
	p.age++
}

func (Person) justLikeThat() {
	fmt.Println("Just like that. No Person instance needed")
}

func main() {
	p := Person{firstName: "John", lastName: "Doe", age: 21, PhoneNumber: PhoneNumber{countryCode: 91, cell: "8689905873"}}

	p1 := Person{age: 30}

	fmt.Printf("Age and name for p: %v and %v\n", p.age, p.firstName)
	fmt.Printf("Age and name for p1: %v and %v\n", p1.age, p1.firstName)

	user := struct {
		username string
		email    string
	}{
		username: "hi@1",
		email:    "whatdahell@1.com",
	}

	fmt.Println("User = ", user)

	fmt.Println("fullname using method: ", p.fullName())

	p.add1ToAge()
	fmt.Println("Printing age after adding 1 to age: ", p.age)
	fmt.Println("Directly accessing cell: ", p.cell)
}
