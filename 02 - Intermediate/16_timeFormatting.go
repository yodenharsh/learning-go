package main

import (
	"fmt"
	"time"
)

func main() {
	layout := "2006-01-02T15:04:05Z07:00"
	str := "2024-07-04T14:30:18Z"

	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println("Encountered error: ", err)
		return
	}

	fmt.Println(t)

	t, err = time.Parse("Jan 02, 2006 03:04 PM", "Jul 03, 2024 03:18 PM")
	fmt.Println(t)
}
