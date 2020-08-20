package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/provider/session"
)

func (a *App) authCheck(r *http.Request) (entity.Session, error) {
	sessionCookie, err := r.Cookie(session.CookieName)
	if err != nil {
		return entity.Session{}, ErrUnauthorized
	}

	token := sessionCookie.Value

	userSession, err := a.Service.GetCurrentSession(token)
	if err != nil {
		if errors.Is(err, hub.ErrMissingValidSession) {
			return entity.Session{}, ErrUnauthorized
		}

		return entity.Session{}, err
	}

	return userSession, nil
}

func createLoginCookie(w http.ResponseWriter, userSession entity.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     session.CookieName,
		Value:    userSession.Token,
		Expires:  time.Now().Add(session.Duration),
		Path:     "/",
		HttpOnly: true,
	})
}

func deleteLoginCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     session.CookieName,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})
}
