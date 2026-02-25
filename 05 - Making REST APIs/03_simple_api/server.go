package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := 3000

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Handling incoming orders")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Handling incoming users")
	})

	fmt.Println("Server is running on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
