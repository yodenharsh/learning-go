package main

import (
	"fmt"
	"regexp"
)

func main() {

	// email
	re := regexp.MustCompile(`[a-zA-Z0-9._+%-]+@[a-zA-Z0-9.-]+\.[a-zA-z]{2,}`)

	email1 := "example@example.com"
	email2 := "invalidEmail"

	fmt.Println("Email1: ", re.MatchString(email1))
	fmt.Println("Email2: ", re.MatchString(email2))

	// 	Capturing group (whatever that is)
	//. Capture date components
	re = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})`)

	date1 := "2024-01-01"
	// date2 := "2025-01-111"

	subMatches := re.FindStringSubmatch(date1)

	fmt.Println(subMatches)
}
