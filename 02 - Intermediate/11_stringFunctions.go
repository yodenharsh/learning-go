package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	str := "Hello everybody"

	fmt.Println(len(str))

	str1 := "Hello there"
	str2 := "But like what"

	result := str1 + " " + str2
	fmt.Println(result)

	fmt.Printf("First element of str1: %c\n", str[0])

	// String conversion

	num := 18
	str3 := strconv.Itoa(num)
	fmt.Println(len(str3))

	// string splitting

	fruits := "apple,orange,banana"
	parts := strings.Split(fruits, ",")

	fmt.Println(parts)

	// joining

	countries := []string{"India", "Pakistan", "Nepal", "USA"}
	countriesString := strings.Join(countries, ",")

	fmt.Println(countriesString)

	// checking if includes

	doesContainInd := strings.Contains(countriesString, "Ind")
	fmt.Println("Does contain ind:", doesContainInd)

	// replacing stuff

	replaced := strings.Replace(countriesString, "a", "Replaced", -1)
	fmt.Println("Replaced: ", replaced)

	// repeating?

	fmt.Println(strings.Repeat(fruits, 3))

	// Counting character
	fmt.Println(strings.Count("Helloo", "o"))

	// RegEx
	re := regexp.MustCompile(`[0-9]+`)
	fmt.Println("true regex: ", re.FindAllString("Hell 12 go878", -1))
	fmt.Println("false regex: ", re.MatchString("Hello"))

	// Strings.builder
	var builder strings.Builder

	builder.WriteString("Hello")
	builder.WriteString(" World")

	resString := builder.String()
	fmt.Println("Built string: ", resString)
}
