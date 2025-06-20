package main

import (
	"fmt"
	"strconv"
)

func main() {
	numberString := "123"
	num, err := strconv.Atoi(numberString)
	if err != nil {
		fmt.Println("Error parsing the value: ", err)
	}

	fmt.Println("Parsed integer + 1: ", num+1)
}
