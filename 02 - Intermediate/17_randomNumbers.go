package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(rand.Intn(101))

	randWithSeed := rand.New(rand.NewSource(42))
	fmt.Println(randWithSeed.Intn(2000))

	for {
		fmt.Println("1. Roll the dice")
		fmt.Println("2. Exit")

		var choice int
		fmt.Scanf("%d", &choice)

		if choice == 1 {
			fmt.Println("You rolled: ", randWithSeed.Intn(6)+1)
			fmt.Scanf("%d", &choice)
		} else {
			fmt.Println("Exitting")
			break
		}
	}
}
