package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type Person struct {
	FirstName string  `json:"name"`
	Age       int     `json:"age,omitempty"`
	Email     string  `json:"email"`
	Address   Address `json:"address"`
}

func main() {
	person := Person{FirstName: "John", Email: "harshmorayya@ymail.com"}
	jsonData, err := json.Marshal(person)

	if err != nil {
		fmt.Println("Could not marshal person")
		return
	}

	fmt.Println(string(jsonData))

	person.Address = Address{City: "New Bombay", State: "MH"}

	jsonData, _ = json.Marshal(person)

	fmt.Println("After adding address: ", string(jsonData))

	rawJsonString := `{"name": "Someone", "age": 12, "email": "harshmorayya3@gmail.com","address": {"city": "San Jose", "state":"CA"}}`
	var personFromJson Person
	err = json.Unmarshal([]byte(rawJsonString), &personFromJson)

	if err != nil {
		fmt.Println("Could not unmarshal: ", err)
		return
	}

	fmt.Println(personFromJson)

	// Handling unknown structures
	jsonData2 := `{"name": "john", "age": 30, "address" : {"city": "something", "state": "NY"}, "email": "2@3.com"}`
	var data map[string]any

	if err = json.Unmarshal([]byte(jsonData2), &data); err != nil {
		fmt.Println("Error unmarshalling: ", err)
		return
	}

	fmt.Println("Unmarshalled JSON: ", data)
	fmt.Println("Accessing specific property: ", data["address"])
}
