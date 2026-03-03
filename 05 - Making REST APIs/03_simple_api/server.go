package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

type PersonInfo struct {
	Age     int    `json:"age"`
	Country string `json:"country"`
}

func main() {
	port := 3000

	people := map[string]PersonInfo{
		"alice":     {Age: 30, Country: "USA"},
		"bob":       {Age: 25, Country: "UK"},
		"Curry boy": {Age: 35, Country: "India"},
	}

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		logRequestDetails(r)
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Handling incoming orders")
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		logRequestDetails(r)
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			http.Error(w, "Missing 'name' query parameter", http.StatusBadRequest)
			return
		}

		person, exists := people[name]
		if !exists {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(person); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	})

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port", port)

	server.ListenAndServeTLS(cert, key)
}

func logRequestDetails(r *http.Request) {
	httpVersion := r.Proto
	fmt.Println("Received request with HTTP version: ", httpVersion)

	tlsVersion := r.TLS.Version
	fmt.Println("TLS version used: ", tlsVersion)
}
