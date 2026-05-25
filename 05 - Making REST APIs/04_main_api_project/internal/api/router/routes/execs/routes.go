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
	mux.HandleFunc("POST /execs/forgotpassword", handlers.ForgotPasswordHandler)
	mux.HandleFunc("POST /execs/password/reset/{resetCode}", handlers.ResetPasswordHandler)
	mux.Handle("POST /execs/{id}/updatepassword", mw.JwtMiddleware(http.HandlerFunc(handlers.UpdatePasswordHandler)))
}
