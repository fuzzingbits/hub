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

// Route Names
const (
	RouteServerStatus = "/api/server/status"
	RouteServerSetup  = "/api/server/setup"
	RouteUserMe       = "/api/user/me"
	RouteUserLogin    = "/api/user/login"
)

// GetRoutes gets all the routes
func (a *App) GetRoutes() []rooter.Route {
	return []rooter.Route{
		{
			Path:     RouteServerStatus,
			Handler:  rooter.ResponseFunc(a.handlerServerStatus),
			Response: entity.ServerStatus{},
		},
		{
			Path:     RouteServerSetup,
			Handler:  rooter.ResponseFunc(a.handlerServerSetup),
			Payload:  entity.CreateUserRequest{},
			Response: entity.UserContext{},
		},
		{
			Path:     RouteUserLogin,
			Handler:  rooter.ResponseFunc(a.handlerUserLogin),
			Response: entity.UserContext{},
			Payload:  entity.UserLoginRequest{},
		},
		{
			Path:     RouteUserMe,
			Handler:  rooter.ResponseFunc(a.handlerUserMe),
			Response: entity.UserContext{},
			Middleware: []rooter.Middleware{
				a.middlewareRequireAuth,
			},
		},
	}
}

// RegisterRoutes for the API
func RegisterRoutes(mux *http.ServeMux, service *hub.Service) {
	a := &App{
		Service: service,
	}

	rooter.RegisterRoutes(mux, a.GetRoutes(), []rooter.Middleware{
		a.middlewareLogger,
		a.middlewareRecovery,
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
		return rooter.ResponseBadRequest()
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
	sessionCookie, _ := req.Cookie(session.CookieName)
	token := sessionCookie.Value

	userSession, _ := a.Service.GetCurrentSession(token)

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userSession.Context,
	}
}

func (a *App) handlerUserLogin(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Build the login request
	var loginRequest entity.UserLoginRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		return rooter.ResponseBadRequest()
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
