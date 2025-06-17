package main

import (
	"errors"
	"fmt"
	"math"
)

type customErr struct {
	message string
}

func (err *customErr) Error() string {
	return fmt.Sprintf("Error: %s", err.message)
}

func eprocess() error {
	return &customErr{message: "Anonymous error created"}
}

func sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, errors.New("Math Error: square root of negative number not possible")
	}

	return math.Sqrt(x), nil
}

func processData(arr []byte) error {
	if len(arr) == 0 {
		return errors.New("length is 0")
	} else {
		return nil
	}
}

func main() {
	result, err := sqrt(16)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Response: ", result)
	data := []byte{}
	if processDataErr := processData(data); processDataErr != nil {
		fmt.Println("Error on processData")
	}

	if eprocessErr := eprocess(); eprocessErr != nil {
		fmt.Println("Error was there on eprocess")
	}

}
