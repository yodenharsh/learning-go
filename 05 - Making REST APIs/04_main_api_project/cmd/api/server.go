package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Teacher struct {
	Id        int    `json:"id,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Class     string `json:"class,omitempty"`
	Subject   string `json:"subject,omitempty"`
}

var (
	teachers = make(map[int]Teacher)
	nextId   = 1
	mutex    = &sync.Mutex{}
)

func init() {
	teachers[nextId] = Teacher{
		Id:        nextId,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Math",
	}
	nextId++
	teachers[nextId] = Teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10B",
		Subject:   "Science",
	}
	nextId++
	teachers[nextId] = Teacher{
		Id:        nextId,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "4C",
		Subject:   "English",
	}
	nextId++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimPrefix(path, "/")

	if id == "" {
		firstName := r.URL.Query().Get("firstName")
		lastName := r.URL.Query().Get("lastName")

		teacherList := make([]Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if firstName != "" && teacher.FirstName != firstName {
				continue
			} else if lastName != "" && teacher.LastName != lastName {
				continue
			}
			teacherList = append(teacherList, teacher)
		}

		response := struct {
			Status string    `json:"status"`
			Count  int       `json:"count"`
			Data   []Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		id, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		teachers, exists := teachers[id]
		if !exists {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(teachers)
	}
}

func postTeachersHandle(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)

	if err != nil {
		http.Error(w, "error decoding JSON", http.StatusUnprocessableEntity)
		return
	}

	addedTeachers := make([]Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		newTeacher.Id = nextId
		teachers[nextId] = newTeacher
		addedTeachers[i] = newTeacher
		nextId++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}

	json.NewEncoder(w).Encode(response)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello root route"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		postTeachersHandle(w, r)
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
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class", "firstName", "lastName"},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", execsHandler)

	rl := mw.NewRateLimiter(5, time.Minute)

	// Middlewares order is first-in first-applied
	secureMux := applyMiddlewares(mux,
		mw.Hpp(hppOptions),
		mw.Compression,
		mw.SecurityHeaders,
		mw.ResponseTimeMiddleware,
		rl.RateLimitingMiddleware,
		mw.Cors)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: secureMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func applyMiddlewares(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
