package main

import "fmt"

func main() {
	var a int = 32;
	b := int32(a)
	var _ = float64(b)
	
	e := 3.12
	f := int(e)
	fmt.Println("Converting 3.12 to f: ", f)

	// Converting string to byte slice
	aString := "Hi there :)"
	var h []byte
	h = []byte(aString)
	fmt.Println(h)
	// Converting byte slice to string
	i := []byte{255,12,71}
	j := string(i)
	fmt.Println(j)
}