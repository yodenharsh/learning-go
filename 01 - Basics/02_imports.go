package main

import (
	"fmt"
	foo "net/http"
)

func main() {
	fmt.Println("Hello everybody")

	resp, err := foo.Get("https://jsonplaceholder.typicode.com/posts/1")

	if err != nil {
		fmt.Println("Error occured:", err)
		return
	}

	defer resp.Body.Close()

	fmt.Println(resp.Status)
}
