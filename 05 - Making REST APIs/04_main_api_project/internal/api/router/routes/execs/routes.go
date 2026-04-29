package execs

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs", handlers.ExecsHandler)
	mux.HandleFunc("POST /execs", handlers.PostExecsHandler)
}
