package main

import (
	"errors"
	"fmt"
)

func SomethingPublic(y, z int) int {
	var x = 1

	return x
}

func somethingPrivate(pt *int) int {
	*pt += 2

	var x = func() {
		fmt.Println("Dobble")
	}
	x()

	return *pt
}

func applyOperation(x, y int, operation func(int, int) int) any {

	return operation(x, y)
}

func multipleReturnValues() (int, int) {
	return 1, 2
}

func propogateErr(a, b int) (string, error) {
	if a > b {
		return "ok", nil
	} else {
		return "", errors.New("Lmao it failed")
	}
}

func variadicArgFunction(x int, optionalInfiniteX ...string) {
	fmt.Println(optionalInfiniteX)
	defer fmt.Println("INFINITE RECURSION")

	variadicArgFunction(1, []string{"hello", "pwe", "god"}...)
}
