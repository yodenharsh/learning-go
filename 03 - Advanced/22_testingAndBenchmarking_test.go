package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func GenerateRandomSlice(size int) []int {
	slice := make([]int, size)
	for i := range slice {
		slice[i] = rand.Int()
	}
	return slice
}

func sumSlice(slice []int) int {
	sum := 0
	for _, num := range slice {
		sum += num
	}
	return sum
}

func TestGenerateRandomSlice(t *testing.T) {
	size := 100
	slice := GenerateRandomSlice(size)
	if len(slice) != size {
		t.Errorf("Expected slice of size %d, got %d", size, len(slice))
	}
}

func BenchmarkGenerateRandomSlice(b *testing.B) {
	for b.Loop() {
		GenerateRandomSlice(1000)
	}
}

func BenchmarkSumSlice(b *testing.B) {
	slice := GenerateRandomSlice(1000)
	// b.Loop() internally calls b.ResetTimer() before the loop and b.StopTimer() after the loop, so
	// we don't need to.
	for b.Loop() {
		sumSlice(slice)
	}
}

func Add(a, b int) int {
	return a + b
}

// Benchmarking

func BenchmarkAddSmallInput(b *testing.B) {
	for b.Loop() {
		Add(2, 3)
	}
}

func BenchmarkAddLargeInput(b *testing.B) {
	for b.Loop() {
		Add(2e18, 3e18)
	}

}

// Testing

func TestAddSubtests(t *testing.T) {
	tests := []struct{ a, b, expected int }{
		{2, 3, 5},
		{0, 0, 0},
		{-1, 1, 0},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("Add(%d, %d)", test.a, test.b), func(t *testing.T) {
			result := Add(test.a, test.b)
			if result != test.expected {
				t.Errorf("Add(%d, %d) = %d; want %d", test.a, test.b, result, test.expected)
			}
		})
	}
}

func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		a, b, expected int
	}{
		{2, 3, 5},
		{0, 0, 0},
		{-1, 1, 0},
	}

	for _, test := range tests {
		result := Add(test.a, test.b)
		if result != test.expected {
			t.Errorf("Add(%d, %d) = %d; want %d", test.a, test.b, result, test.expected)
		}
	}
}

func TestAdd(t *testing.T) {
	result := Add(2, 3)
	if result != 5 {
		t.Errorf("Add(2, 3) = %d; want 5", result)
	}
}
