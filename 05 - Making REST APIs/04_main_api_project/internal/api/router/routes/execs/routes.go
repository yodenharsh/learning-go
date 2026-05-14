package execs

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("GET /execs/{id}", handlers.GetExecByIdHandler)
	mux.HandleFunc("POST /execs", handlers.PostExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchExecByIdHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecByIdHandler)

	// GetExecsHandler is temporary. Actual handlers need to be implemented for these routes
	mux.HandleFunc("POST /execs/login", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/logout", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/forgotpassword", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/password/reset/{resetCode}", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.GetExecsHandler)
}
