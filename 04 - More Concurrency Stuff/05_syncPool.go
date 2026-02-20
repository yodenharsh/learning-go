package main

import (
	"fmt"
	"sync"
)

type person struct {
	name string
	age  int
}

func main() {
	var pool = sync.Pool{
		New: func() any {
			fmt.Println("Creating new person")
			return &person{}
		},
	}

	person1 := pool.Get().(*person)
	person1.name = "John"
	person1.age = 20
	fmt.Printf("Person 1: %+v\n", person1)

	pool.Put(person1)
	fmt.Println("Returned person1 to pool")

	person2 := pool.Get().(*person)
	fmt.Println("Got person2 from pool: ", person2)

	person3 := pool.Get()
	if person3 != nil {
		fmt.Println("Got person3 from pool: ", person3)
	} else {
		fmt.Println("Sync pool is empty")
	}

	// Returning all other objects back to the pool
	pool.Put(person2)
	pool.Put(person3)
	fmt.Println("Returned person2 and person3 to pool")

	person4 := pool.Get().(*person)
	fmt.Println("Got person4 from pool: ", person4)

	person5 := pool.Get().(*person)
	fmt.Println("Got person5 from pool: ", person5)
}
