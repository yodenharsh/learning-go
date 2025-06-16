package main

import (
	"fmt"
	"math"
)

type geometry interface {
	area() float64
	perim() float64
}

type rect struct {
	width, height float64
}

func (r rect) area() float64 {
	return r.height * r.width
}
func (r rect) perim() float64 {
	return 2 * (r.height + r.width)
}

type circle struct {
	radius float64
}

func (c circle) area() float64 {
	return math.Pi * math.Pow(c.radius, 2)
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}
func (c circle) diameter() float64 {
	return 2 * c.radius
}

func measure(g geometry) {
	fmt.Println("The geometry: ", g)
	fmt.Println("The area: ", g.area())
	fmt.Println("The perimeter: ", g.perim())
}

func main() {
	r := rect{width: 44, height: 3}
	c := circle{radius: 21}

	measure(r)
	measure(c)
}

func checkType(g geometry) {
	switch g.(type) {
	case circle:
		fmt.Println("It's a circle")
	case rect:
		fmt.Println("It's a rectange")
	default:
		fmt.Println("How did we even get here")
	}
}
