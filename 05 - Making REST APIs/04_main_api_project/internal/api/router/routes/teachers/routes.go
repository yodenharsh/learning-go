package teachers

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherByIdHandler)
	mux.HandleFunc("POST /teachers", handlers.PostTeachersHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeachersHandler)
	mux.HandleFunc("PATCH /teachers", handlers.PatchTeachersHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherByIdHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)
}
