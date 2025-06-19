package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println(time.Now())

	specificTime := time.Date(2024, time.July, 12, 0, 0, 0, 0, time.UTC)
	fmt.Println("Specific time: ", specificTime)

	// See https://pkg.go.dev/time#Layout
	parsedTime, _ := time.Parse("2006-01-02", "2020-05-01")

	fmt.Println(parsedTime)

	// Formatting time instead of string now
	t := time.Now()
	fmt.Println("Formatted time: ", t.Format("2006/02/01 Mon"))

	oneDayLater := t.Add(time.Hour * 24)
	fmt.Println(oneDayLater)

	loc, _ := time.LoadLocation("Asia/Kolkata")
	t = time.Date(2024, time.July, 8, 14, 16, 40, 00, time.UTC)
	tLocal := t.In(loc)
	fmt.Println(tLocal)

	loc, _ = time.LoadLocation("America/New_York")

	// Convert time to location
	tInNy := time.Now().In(loc)
	fmt.Println(tInNy)
}
