package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

type ByAge []Person

func (a ByAge) Len() int { return len(a) }
func (a ByAge) Less(i, j int) bool {
	return a[i].Age < a[j].Age
}
func (a ByAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func main() {
	numbers := []int{5, 2, 9, 1, 5, 6}
	sort.Ints(numbers) // Inplace operation
	fmt.Println("Sorted numbers:", numbers)

	stringSlice := []string{"banana", "apple", "cherry"}
	sort.Strings(stringSlice) // Inplace operation
	fmt.Println("Sorted strings:", stringSlice)

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}
	sort.Sort(ByAge(people)) // Also an inplace operation
	fmt.Println("Sorted people by age:", people)

	// Sort by slice here
	stringSlice = []string{"banana", "apple", "cherry"}
	sort.Slice(stringSlice, func(i, j int) bool {
		return len(stringSlice[i]) < len(stringSlice[j])
	})
	fmt.Println("Sorted strings using sort.Slice:", stringSlice)

}
