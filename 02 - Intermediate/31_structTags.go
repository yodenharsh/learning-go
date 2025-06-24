package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	FirstName string `json:"-"` // Firstname should always be omitted when "-" as key
	LastName  string `json:"lastName,omitempty"`
	Age       int    `json:"age"`
}

func main() {
	person := Person{
		FirstName: "Jane",
		LastName:  "", // This is still omitted when omitempty
		Age:       0,  // THis is also omitted when omitempty is used
	}

	jsonData, err := json.Marshal(person)
	if err != nil {
		log.Fatalln("Error occurred when unmarshalling: ", err)
		return
	}

	fmt.Println(string(jsonData))
}
