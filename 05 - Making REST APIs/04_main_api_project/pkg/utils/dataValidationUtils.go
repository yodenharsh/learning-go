package utils

import (
	"errors"
	"fmt"
	"reflect"
)

func CheckStringFieldsNotEmpty(value reflect.Value) error {
	val := reflect.ValueOf(value)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Kind() == reflect.String && field.String() == "" {
			fmt.Println("Field name: ", field.Type().Name())
			return ErrorHandler(errors.New("One or more required fields is empty"), "One or more required fields is empty", 400)
		}
	}
	return nil
}
