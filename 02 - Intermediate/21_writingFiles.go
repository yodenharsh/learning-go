package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("output.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data := []byte("Hello World!\n")
	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file: ", err)
		return
	}

	fmt.Println("Data has been written successfully")

}
