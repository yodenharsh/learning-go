package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	contextTodoAndBkg()

	ctx := context.TODO()
	result := isEvenOdd(ctx, 10)
	fmt.Println(result)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Cancel is being called manually but will also be called automatically when the timeout expires
	// Note: context.withCancel() is also an option if you want to cancel the context manually without a timeout
	result = isEvenOdd(ctx, 5)
	fmt.Println("Result from context with timeout:", result)

	time.Sleep(3 * time.Second)
	result = isEvenOdd(ctx, 5)
	fmt.Println("Result from context with timeout after sleeping:", result)

	rootCtx := context.Background()
	ctx, cancel = context.WithTimeout(rootCtx, 2*time.Second)
	defer cancel()

	ctx = context.WithValue(ctx, "requestId", "124jl")

	go doWork(ctx)
	time.Sleep(3 * time.Second)
	requestId := ctx.Value("requestId")
	if requestId != nil {
		fmt.Println("Request ID:", requestId)
	} else {
		fmt.Println("Request ID not found")
	}
}

func contextTodoAndBkg() {
	todoContext := context.TODO()      // Used when you are not sure which context to use or when the context is not yet available
	contextBkg := context.Background() // Usually used as the root context for other contexts

	ctx := context.WithValue(todoContext, "name", "John")
	fmt.Println(ctx)
	fmt.Println(ctx.Value("name"))

	ctx1 := context.WithValue(contextBkg, "city", "New York")
	fmt.Println(ctx1)
	fmt.Println(ctx1.Value("city"))
}

func isEvenOdd(ctx context.Context, num int) string {
	select {
	case <-ctx.Done(): // Check docs of Done() by hovering over it
		return "Cancelled"
	default:
		if num%2 == 0 {
			return "Even"
		}
		return "Odd"
	}
}

func doWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Work cancelled: ", ctx.Err())
			return
		default:
			fmt.Println("Working...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}
