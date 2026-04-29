package root

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", handlers.RootHandler)
}
