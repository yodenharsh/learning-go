package middlewares

import (
	"log"
	"net/http"
	"time"
)

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	wrappedWriter := &responseWriter{
		ResponseWriter: w,
		status: http.StatusOK,
	}

	next.ServeHTTP(wrappedWriter,r)
	duration := time.Since(start)
	// Assume we are logging
	log.Printf("%s %s %d\n%v\n\n", r.Method, r.URL.Path, wrappedWriter.status, duration)
})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}