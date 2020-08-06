package api

import (
	"encoding/json"
	"net/http"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
)

// Route Names
const (
	RouteServerStatus = "/api/server/status"
	RouteServerSetup  = "/api/server/setup"
	RouteUserMe       = "/api/user/me"
	RouteUserLogin    = "/api/user/login"
)

// App for the REST API
type App struct {
	Service *hub.Service
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
		},
	}
}

func (a *App) handlerServerStatus(w http.ResponseWriter, req *http.Request) rooter.Response {
	serverStatus, err := a.Service.GetServerStatus()
	if err != nil {
		return a.generateErrorResponse(err, req)
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
		return rooter.ResponseBadRequest
	}

	userSession, err := a.Service.SetupServer(payload)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	createLoginCookie(w, userSession)

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userSession.Context,
	}
}

func (a *App) handlerUserMe(w http.ResponseWriter, req *http.Request) rooter.Response {
	userSession, err := a.authCheck(req)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

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
		return rooter.ResponseBadRequest
	}

	userSession, err := a.Service.Login(loginRequest)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	createLoginCookie(w, userSession)

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userSession.Context,
	}
}
