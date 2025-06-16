package main

import "fmt"

func swap[T comparable](a, b *T) {
	var temp T = *a
	*a = *b
	*b = temp
}

type stack[T any] struct {
	elements []T
}

func (s *stack[T]) push(el T) {
	s.elements = append(s.elements, el)
}
func (s *stack[T]) pop() (T, bool) {
	if len(s.elements) == 0 {
		var zero T
		return zero, false
	}

	lastEl := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return lastEl, true
}
func (s stack[T]) isEmpty() bool {
	return len(s.elements) == 0
}
func (s stack[T]) printAll() {
	if s.isEmpty() {
		fmt.Println("No elements in stack")
		return
	}

	for _, v := range s.elements {
		fmt.Print("%v ", v)
	}
	fmt.Println()
}

func main() {
	a, b := 5, 10

	fmt.Printf("Before swapping: %v and %v\n", a, b)
	fmt.Printf("After swapping: %v and %v", b, a)

	intStack := stack[int]{
		elements: []int{1, 4, 5},
	}
	stringStack := stack[string]{
		elements: []string{"Hello", "there", "troy"},
	}

	intStack.push(4)
	intStack.push(12)

	intStack.printAll()

	poppedEl, wasSuccess := intStack.pop()
	if wasSuccess {
		fmt.Println("Pop the stack: ", poppedEl)
	}

	if stringStack.isEmpty() {
		fmt.Println("stringStack is empty")
	}
}
