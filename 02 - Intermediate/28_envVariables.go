package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	user := os.Getenv("USER")
	home := os.Getenv("HOME")

	fmt.Println("User env: ", user)
	fmt.Println("Home env: ", home)

	for _, e := range os.Environ() {
		kvpair := strings.SplitN(e, "-", 2)
		fmt.Println(kvpair[0])
	}

}
