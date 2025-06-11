package main

import "fmt"

func main() {
	var ptr *int

	var a int = 10
	ptr = &a

	fmt.Println(a)
	fmt.Println(ptr)

	modifyValue(ptr)
	fmt.Println(*ptr)
}

func modifyValue(p *int) {
	*p = *p + 20
}
