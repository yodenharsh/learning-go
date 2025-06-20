package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(strings.NewReader("A fairly large string I think?"))

	// reading the data byte by byte
	data := make([]byte, 20)
	bytesRead, err := reader.Read(data)

	if err != nil {
		fmt.Println("Error reading: ", err)
		return
	}

	fmt.Printf("Number of bytes read: %d. Where bytes: %s\n", bytesRead, data[:bytesRead])

	line, err := reader.ReadString('?')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Line read: ", line)

	writer := bufio.NewWriter(os.Stdout)

	// Writing byte slice
	data = []byte("Hello, bufio package!\n")
	n, err := writer.Write(data)

	if err != nil {
		fmt.Println("Error writing: ", err)
	}

	fmt.Printf("Wrote %d bytes\n", n)
	if err = writer.Flush(); err != nil {
		fmt.Println("Error occurred when flushing: ", err)
	}

	str := "This is a string\n"
	n, err = writer.WriteString(str)
	if err != nil {
		fmt.Println("Error writing string: ", err)
	}
	fmt.Printf("Wrote %d bytes \n", n)
	writer.Flush()
}
