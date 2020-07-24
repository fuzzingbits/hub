package rooter

import (
	"log"
	"net/http"
)

// Middleware for wrapping http.Handlers
type Middleware func(next http.Handler) http.Handler

// MiddlewareLogger is a simple logging middleware
func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		log.Printf("Hey, we logged a request for %s", r.URL.Path)
	})
}
