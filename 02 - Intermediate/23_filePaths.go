package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func main() {

	joinedPath := filepath.Join("downloads", "file.zip")
	fmt.Println("Joined path: ", joinedPath)

	normalizedPath := filepath.Clean("./data/../data/file.txt")
	fmt.Println("Cleaned and normalized path: ", normalizedPath)

	dir, file := filepath.Split("/home/mharsh-ariqt/docs/file.txt")
	fmt.Println("file name: ", file)
	fmt.Println("Directory: ", dir)

	fmt.Println(strings.TrimSuffix(file, filepath.Ext(file)))
}
