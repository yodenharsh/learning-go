package main

import (
	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
	"time"
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

	switch r.Method {
	case http.MethodGet:

		w.Write([]byte("Hello GET Method on Teachers Route"))
	case http.MethodPost:
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

	if r.Method == http.MethodPost {
		w.Write([]byte("Hello POST Method on Execs Route"))
		err := r.ParseForm()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wronge when trying to access form values", http.StatusInternalServerError)
			return
		}
		formMap := r.Form
		fmt.Println("formMap:", formMap)

		fmt.Println("Queries: ", r.URL.Query())
	}
}

func main() {

	port := "3000"
	fmt.Println("Server is running on port ", port)
	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	rl := mw.NewRateLimiter(5, time.Minute)
	secureMux := rl.RateLimitingMiddleware(
		mw.ResponseTimeMiddleware(
			mw.Compression(
				mw.SecurityHeaders(
					mw.Cors(mux)))))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mw.Hpp(hppOptions)(secureMux),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
