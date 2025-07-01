package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func readFromReader(r io.Reader) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		log.Fatalln("Error while reading buffer: ", err)
	}
	fmt.Println("Read ",n," bytes")
	fmt.Println(string(buf[:n]))
}

func writeToWriter(w io.Writer, data string) {
	_, err := w.Write([]byte(data))
	if err != nil {
		log.Fatalln("Error while writing data to a buffer: ", err)
	}
}

func closeResource(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatalln("Error while closing: ", err)
	} 
}

func bufferExample() {
	var buf bytes.Buffer
	buf.WriteString("Hello buffer!")
	fmt.Println(buf.String())
}

func multiReaderExample() {
	r1 := strings.NewReader("Hello")
	r2 := strings.NewReader("World!")
	mr := io.MultiReader(r1,r2)

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(mr);
	if err != nil {
		log.Fatalln("Error occurred: ", err)
	}
	fmt.Println(buf.String())
}

func pipeExample() {
	pr, pw := io.Pipe()
	go func() {
		pw.Write([]byte("Hello pipe"))
		pw.Close()
	}()

  buf := new(bytes.Buffer)
  buf.ReadFrom(pr)
  fmt.Println(buf.String())
}

func writeToFile(filepath string, data string) {
	file, err := os.OpenFile(filepath,os.O_APPEND | os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Error opening/creating file: ", err)
	}
	defer closeResource(file)
	
	// writer := io.Writer(file)
	// if _, err := writer.Write([]byte("Writing to the file")); err != nil {
	// 	log.Fatalln("Error: ", err)
	// }

	if _, err := file.Write([]byte(data)); err != nil {
		log.Fatalln("Error when writing to file: ", err)
	}
}

func main() {
	fmt.Println("Read from reader")
	readFromReader(strings.NewReader("What even is going on")) 

	fmt.Println("Writing to Writer")
	var writer bytes.Buffer
	writeToWriter(&writer,"Hello writer")

	fmt.Println("Simple buffer example")
	bufferExample()

	fmt.Println("Multi reader example")
	multiReaderExample()

	fmt.Println("Pipe example")
	pipeExample()

	filepath := "io.txt"
	writeToFile(filepath, "Adding to io")
	writeToFile(filepath, "Adding another line")
}