package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s\n", r.Method, r.URL.Path)
	fmt.Fprintln(w, "Hello server")
	w.WriteHeader(http.StatusOK)
}

func main() {

	const serverAddr string = ":3000"

	http.HandleFunc("/", handler)
	log.Printf("Starting server at %s\n", serverAddr)

	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		log.Fatalln(err)
	}
}
