package main

import (
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	file, err := os.CreateTemp("", "temporaryFile")

	checkError(err)

	fmt.Println("Wrote a temporary file: ", file.Name())

	defer os.Remove(file.Name())
	defer file.Close()

	tempDir, err := os.MkdirTemp("", "GoTempDir")
	checkError(err)

	defer os.RemoveAll(tempDir)
	fmt.Println("Temporary directory: ", tempDir)
}
