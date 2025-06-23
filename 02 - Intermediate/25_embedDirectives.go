package main

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed example.txt
var content string

//go:embed "01 - Basics"
var basicsFolder embed.FS

func main() {
	fmt.Println("Embedded content: ", content)
	content, err := basicsFolder.ReadFile("01 - Basics/01_hello.go")
	if err != nil {
		fmt.Println("Error reading folder: ", err)
	}

	fmt.Println("Embedded file content: ", string(content))

	fs.WalkDir(basicsFolder, "01 - Basics", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(path)
		return nil
	})
}
