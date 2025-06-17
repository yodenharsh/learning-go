package main

import (
	"errors"
	"fmt"
)

type customError struct {
	code    int
	message string
	err     error
}

func (e *customError) Error() string {
	return fmt.Sprintf("Error %d: %s\n", e.code, e.message)
}

func doSomething() error {
	return &customError{
		code:    500,
		message: "Internal server error",
	}
}

func doSomething1() error {
	err := doSomethingElse()
	if err != nil {
		return &customError{
			code:    500,
			message: "Something went wrong",
			err:     err,
		}
	}

	return nil
}

func doSomethingElse() error {
	return errors.New("Internal error")
}

func main() {
	if err := doSomething(); err != nil {
		fmt.Println("Error occurred: ", err)
	}

	if err := doSomething1(); err != nil {
		fmt.Println("Error occurred: ", err)
	}
}
