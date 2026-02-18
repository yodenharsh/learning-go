package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

type Greeter struct{}

func (g Greeter) Greet(name string) string {
	return "Hello, " + name + "!"
}

func main() {
	x := 42
	v := reflect.ValueOf(x)
	t := v.Type()

	fmt.Println("Value:", v)
	fmt.Println("Type:", t)
	fmt.Println("Kind: ", t.Kind())
	fmt.Println("Is Int: ", t.Kind() == reflect.Int)
	fmt.Println("Is string: ", t.Kind() == reflect.String)
	fmt.Println("Is zero: ", v.IsZero())

	y := 10
	v1 := reflect.ValueOf(&y).Elem()
	v2 := reflect.ValueOf(&y)

	fmt.Println("V2 type: ", v2.Type())

	fmt.Println("Original Value: ", v1.Int()) // v.Int() panics if v is not an int
	v1.SetInt(18)
	fmt.Println("New Value: ", v1.Int())

	var itf any = "Hello, Reflection!"
	v3 := reflect.ValueOf(itf)

	fmt.Println("V3 type: ", v3.Type())
	fmt.Println("Interface kind: ", v3.Kind())

	// Working with structs here

	p := Person{Name: "Alice", Age: 30}
	v4 := reflect.ValueOf(p)

	for i := range v4.NumField() {
		fmt.Printf("Field %d: %v\n", i, v4.Field(i))
	}

	v5 := reflect.ValueOf(&p).Elem()
	nameField := v5.FieldByName("name")             // cannot set "name" field because "name" is unexported (lowercase 'n')
	nameFieldWithCapitalN := v5.FieldByName("Name") // can set "Name" field because "Name" is exported (uppercase 'N')
	if nameField.CanSet() {
		nameField.SetString("Bob")
	} else if nameFieldWithCapitalN.CanSet() {
		nameFieldWithCapitalN.SetString("Bob")
	} else {
		fmt.Println("Cannot set the name field")
	}

	fmt.Printf("Updated person: %+v\n", p)

	// Working with methods here
	v6 := reflect.TypeFor[Greeter]()
	// OR
	// 	g := Greeter{}
	//  v6 := reflect.TypeOf(g)

	var method reflect.Method

	fmt.Println("Type: ", v6)
	for i := range v6.NumMethod() {
		method = v6.Method(i)
		fmt.Printf("Method %d: %s\n", i, method.Name)
	}

	m, wasFound := v6.MethodByName(method.Name)
	if wasFound {
		fmt.Printf("Method found: %s\n", m.Name)
		results := m.Func.Call([]reflect.Value{reflect.ValueOf(Greeter{}), reflect.ValueOf("Alice")})
		fmt.Printf("Method result: %v\n", results[0])
	} else {
		fmt.Println("Method not found")
	}
}
