package api

import "net/http"

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
