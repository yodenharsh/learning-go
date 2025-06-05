package main

import (
	"errors"
	"fmt"
)

func main() {
	process(1)
}

func process(input int) (string, error) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("We were able to handle the panic")
		}
	}()
	defer fmt.Println("process completed")

	if input != 0 {
		panic("Wow")
		return "", errors.New("Wow")
	}
	fmt.Println("We take these")
	return "", nil
}
