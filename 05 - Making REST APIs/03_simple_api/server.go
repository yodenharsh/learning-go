package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"golang.org/x/net/http2"
)

func main() {
	port := 3000

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
		fmt.Fprintf(w, "Handling incoming users")
	})

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: tlsConfig,
		Handler:   nil,
	}

	http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("Server is running on port", port)

	server.ListenAndServeTLS(cert, key)
	// err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	// if err != nil {
	// 	log.Fatalln("Error starting server:", err)
	// }
}

func logRequestDetails(r *http.Request) {
	httpVersion := r.Proto
	fmt.Println("Received request with HTTP version: ", httpVersion)

	tlsVersion := r.TLS.Version
	fmt.Println("TLS version used: ", tlsVersion)
}
