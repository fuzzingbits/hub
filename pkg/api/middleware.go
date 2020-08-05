package api

import (
	"errors"
	"net/http"

	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
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

func (a *App) middlewareRequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie(session.CookieName)
		if err != nil {
			rooter.ResponseUnauthorized().ServeHTTP(w, r)
			return
		}

		token := sessionCookie.Value

		_, err = a.Service.GetCurrentSession(token)
		if err != nil {
			if errors.Is(err, hub.ErrMissingValidSession) {
				rooter.ResponseUnauthorized().ServeHTTP(w, r)
				return
			}

			a.serverError(err, r).ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) middlewareRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err, isErr := recovered.(error)
				if isErr {
					a.serverError(err, r).ServeHTTP(w, r)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
