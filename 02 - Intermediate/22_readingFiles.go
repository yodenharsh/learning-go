package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("output.txt")
	if err != nil {
		fmt.Println("Error opening file")
		return
	}

	defer file.Close()
	fmt.Println("File opened successfully.")

	// data := make([]byte, 1024)
	// _, err = file.Read(data)
	// if err != nil {
	// 	fmt.Println("Error reading data from file: ", err)
	// 	return
	// }

	// fmt.Println("Read file: ", string(data))

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Line: ", line)
	}

	err = scanner.Err()
	if err != nil {
		fmt.Println("Scanner error detected: ", err)
	}
}
