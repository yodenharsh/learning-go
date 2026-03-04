package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := "3000"
	fmt.Println("Server is running on port ", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello root route"))
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Received request for /teachers route on ", r.Method, " method")
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Hello GET Method on Teachers Route"))
		case http.MethodPost:
			w.Write([]byte("Hello POST Method on Teachers Route"))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method not allowed"))
		}
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Placeholder for students route"))
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Placeholder for execs route"))
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
