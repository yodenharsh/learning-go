package main

import (
	"fmt"
)

func main() {

	sampleMap1 := map[string]int{"helllo": 1}

	fmt.Println(sampleMap1["helllo"])

	sampleMap1["sample"] = 5

	delete(sampleMap1, "sample")
	clear(sampleMap1)

	nonExistantValue, isValuePresent := sampleMap1["somekey"]

	if isValuePresent {
		fmt.Printf("nonExistantValue: %v\n", nonExistantValue)
	}

	for key, value := range sampleMap1 {
		fmt.Println(key, value)
	}
}
