package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func main() {
	user := User{Name: "Alice", Email: "alice@example.com"}

	jsonData, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(jsonData))

	var user1 User
	err = json.Unmarshal(jsonData, &user1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User created from JSON data: %+v\n", user1)

	jsonData2 := `{"name": "Bob", "email": "bob@example.com"}`
	reader := strings.NewReader(jsonData2)
	decoder := json.NewDecoder(reader)

	var user2 User
	err = decoder.Decode(&user2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User created from JSON data: %+v\n", user2)

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	err = encoder.Encode(user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User encoded to JSON: %s\n", buf.String())
}
