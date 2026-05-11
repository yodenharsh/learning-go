package students

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /students", handlers.GetStudentsHandler)
	mux.HandleFunc("GET /students/{id}", handlers.GetStudentByIdHandler)
	mux.HandleFunc("POST /students", handlers.PostStudentsHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentsHandler)
	mux.HandleFunc("PATCH /students", handlers.PatchStudentsHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchStudentByIdHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteStudentByIdHandler)
	mux.HandleFunc("DELETE /students", handlers.DeleteStudentsHandler)
}
