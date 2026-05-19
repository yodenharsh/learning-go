package execs

import (
	"net/http"
	"restapi/internal/api/handlers"
	mw "restapi/internal/api/middlewares"
)

func Register(mux *http.ServeMux) {
	mux.Handle("GET /execs", mw.JwtMiddleware(http.HandlerFunc(handlers.GetExecsHandler)))
	mux.Handle("GET /execs/{id}", mw.JwtMiddleware(http.HandlerFunc(handlers.GetExecByIdHandler)))
	mux.Handle("POST /execs", mw.JwtMiddleware(http.HandlerFunc(handlers.PostExecsHandler)))
	mux.Handle("PATCH /execs", mw.JwtMiddleware(http.HandlerFunc(handlers.PatchExecsHandler)))
	mux.Handle("PATCH /execs/{id}", mw.JwtMiddleware(http.HandlerFunc(handlers.PatchExecByIdHandler)))
	mux.Handle("DELETE /execs/{id}", mw.JwtMiddleware(http.HandlerFunc(handlers.DeleteExecByIdHandler)))

	mux.HandleFunc("POST /execs/login", handlers.LoginHandler)
	mux.HandleFunc("POST /execs/logout", handlers.LogoutHandler)
	// GetExecsHandler is temporary. Actual handlers need to be implemented for these routes
	mux.HandleFunc("POST /execs/forgotpassword", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/password/reset/{resetCode}", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", handlers.GetExecsHandler)
}
