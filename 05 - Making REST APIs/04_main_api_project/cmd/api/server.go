package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello root route"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Received request for /teachers route on ", r.Method, " method")
	switch r.Method {
	case http.MethodGet:

		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		userId := strings.TrimSuffix(path, "/")
		fmt.Println("Extracted userId: ", userId)

		w.Write([]byte("Hello GET Method on Teachers Route"))
	case http.MethodPost:
		// Parsing form here
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			log.Println("Error parsing form:", err)
			return
		}

		fmt.Println("Form data received:", r.Form)

		// Prepare response data
		response := make(map[string]string)
		for key, values := range r.Form {
			if len(values) > 0 {
				response[key] = values[0]
			}
		}
		fmt.Println("Response data prepared:", response)

		// Parsing raw json
		jsonData, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		defer r.Body.Close()

		fmt.Println("Raw body: ", string(jsonData))

		var userInstance User
		err = json.Unmarshal(jsonData, &userInstance)
		if err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			log.Println("Error parsing JSON:", err)
			return
		}

		fmt.Println("Parsed JSON data:", userInstance)

		w.Write([]byte("Hello POST Method on Teachers Route"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Placeholder for students route"))
}

func execsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Placeholder for execs route"))
}

func main() {

	port := "3000"
	fmt.Println("Server is running on port ", port)

	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/teachers/", teachersHandler)

	http.HandleFunc("/students/", studentsHandler)

	http.HandleFunc("/execs/", execsHandler)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
