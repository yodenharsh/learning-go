package execs

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs", handlers.ExecsHandler)
	mux.HandleFunc("GET /execs/{id}", handler.)
	mux.HandleFunc("POST /execs", handlers.PostExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.)
	mux.HandleFunc("PATCH /execs/{id}", handlers.)
	mux.HandleFunc("DELETE /execs/{id}", handlers.)

	mux.HandleFunc("POST /execs/login", handlers.)
	mux.HandleFunc("POST /execs/logout", handlers.)
	mux.HandleFunc("POST /execs/forgotpassword", handlers.)
	mux.HandleFunc("POST /execs/password/reset/{resetCode}", handlers.)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.)
}
``