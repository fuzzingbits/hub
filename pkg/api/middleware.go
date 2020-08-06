package api

import (
	"net/http"
)

func (a *App) middlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		a.Service.HTTPLogger.Printf(
			"%s %s",
			r.Method,
			r.URL.Path,
		)
	})
}

func (a *App) middlewareRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err, isErr := recovered.(error)
				if isErr {
					a.generateErrorResponse(err, r).ServeHTTP(w, r)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
