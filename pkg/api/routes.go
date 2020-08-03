package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/rollbar/rollbar-go"
)

// App for the REST API
type App struct {
	Service *hub.Service
}

var responseServerAlreadySetup = rooter.Response{
	StatusCode: http.StatusOK,
	Message:    "Server Is Already Setup",
	State:      false,
}

var responseMissingValidSession = rooter.Response{
	StatusCode: http.StatusOK,
	Message:    "Missing Valid Session",
	State:      false,
}

var responseInvalidLogin = rooter.Response{
	StatusCode: http.StatusOK,
	State:      false,
	Message:    "Invlaid Login",
}

// RegisterRoutes for the API
func RegisterRoutes(mux *http.ServeMux, service *hub.Service) {
	a := &App{
		Service: service,
	}

	routes := []rooter.Route{
		{
			Path:    "/api/server/status",
			Handler: rooter.ResponseFunc(a.handlerServerStatus),
		},
		{
			Path:    "/api/server/setup",
			Handler: rooter.ResponseFunc(a.handlerServerSetup),
		},
		{
			Path:    "/api/user/me",
			Handler: rooter.ResponseFunc(a.handlerUserMe),
		},
		{
			Path:    "/api/user/login",
			Handler: rooter.ResponseFunc(a.handlerUserLogin),
		},
	}

	rooter.RegisterRoutes(mux, routes, []rooter.Middleware{
		a.middlewareLogger,
	})
}

func (a *App) serverError(err error, r *http.Request) rooter.Response {
	// Send to rollbar if Rollbar is configured
	if a.Service.Rollbar != nil {
		a.Service.Rollbar.RequestError(rollbar.ERR, r, err)
	}

	// Print the error to the ErrorLogger
	a.Service.ErrorLogger.Printf(
		"Request Error: %s %s - Err: %s",
		r.Method,
		r.URL.Path,
		err.Error(),
	)

	return rooter.ResponseInternalServerError()
}

func (a *App) handlerServerStatus(w http.ResponseWriter, req *http.Request) rooter.Response {
	serverStatus, err := a.Service.GetServerStatus()
	if err != nil {
		return a.serverError(err, req)
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       serverStatus,
	}
}

func (a *App) handlerServerSetup(w http.ResponseWriter, req *http.Request) rooter.Response {
	var payload entity.CreateUserRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		return a.serverError(err, req)
	}

	userSession, err := a.Service.SetupServer(payload)
	if err != nil {
		if errors.Is(err, hub.ErrServerAlreadySetup) {
			return responseServerAlreadySetup
		}

		return a.serverError(err, req)
	}

	createLoginCookie(w, userSession)

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userSession.Context,
	}
}

func (a *App) handlerUserMe(w http.ResponseWriter, req *http.Request) rooter.Response {
	sessionCookie, err := req.Cookie(session.CookieName)
	if err != nil {
		return responseMissingValidSession
	}

	token := sessionCookie.Value

	session, err := a.Service.GetCurrentSession(token)
	if err != nil {
		if errors.Is(err, hub.ErrMissingValidSession) {
			return responseMissingValidSession
		}

		return a.serverError(err, req)
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       session.Context,
	}
}

func (a *App) handlerUserLogin(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Build the login request
	var loginRequest entity.UserLoginRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		return a.serverError(err, req)
	}

	userSession, err := a.Service.Login(loginRequest)
	if err != nil {
		if errors.Is(err, hub.ErrInvalidLogin) {
			return responseInvalidLogin
		}

		return a.serverError(err, req)
	}

	createLoginCookie(w, userSession)

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userSession.Context,
	}
}

func createLoginCookie(w http.ResponseWriter, userSession entity.Session) {
	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    session.CookieName,
		Value:   userSession.Token,
		Expires: time.Now().Add(session.Duration),
		Path:    "/",
	})
}
