package api

import (
	"encoding/json"
	"errors"
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
	RouteUserLogout   = "/api/user/logout"
	RouteUserNew      = "/api/user/new"
	RouteUserList     = "/api/user/list"
	RouteUserDelete   = "/api/user/delete"
	RouteUserUpdate   = "/api/user/update"
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
			Path:     RouteUserNew,
			Handler:  rooter.ResponseFunc(a.handlerUserNew),
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
			Path:    RouteUserLogout,
			Handler: rooter.ResponseFunc(a.handlerUserLogout),
		},
		{
			Path:     RouteUserMe,
			Handler:  rooter.ResponseFunc(a.handlerUserMe),
			Response: entity.UserContext{},
		},
		{
			Path:     RouteUserList,
			Handler:  rooter.ResponseFunc(a.handlerUserList),
			Response: []entity.User{},
		},
		{
			Path:     RouteUserDelete,
			Handler:  rooter.ResponseFunc(a.handlerUserDelete),
			Response: nil,
		},
		{
			Path:     RouteUserUpdate,
			Handler:  rooter.ResponseFunc(a.handlerUserUpdate),
			Payload:  entity.UpdateUserRequest{},
			Response: entity.UserContext{},
		},
	}
}

func (a *App) handlerUserUpdate(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Require login
	userSession, err := a.authCheck(req)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	// Get the payload
	var payload entity.UpdateUserRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		return rooter.ResponseBadRequest
	}

	userContext, err := a.Service.UpdateUser(payload)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	if userSession.Context.User.UUID == userContext.User.UUID {
		if _, err := a.Service.SaveSession(userSession.Token, userContext); err != nil {
			return a.generateErrorResponse(err, req)
		}
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Message:    "User Updated Successfully",
		Data:       userContext,
	}
}

func (a *App) handlerUserDelete(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Require login
	_, err := a.authCheck(req)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	// Get the payload
	var payload entity.DeleteUserRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		return rooter.ResponseBadRequest
	}

	if err := a.Service.DeleteUser(payload.UUID); err != nil {
		return a.generateErrorResponse(err, req)
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
	}
}

func (a *App) handlerUserList(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Require login
	_, err := a.authCheck(req)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	users, err := a.Service.ListUsers()
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       users,
	}
}

func (a *App) handlerUserNew(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Require login
	_, err := a.authCheck(req)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	// Get the payload
	var payload entity.CreateUserRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&payload); err != nil {
		return rooter.ResponseBadRequest
	}

	userContext, err := a.Service.CreateUser(payload)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userContext,
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
		if errors.Is(err, ErrUnauthorized) {
			deleteLoginCookie(w)
			return rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
			}
		}

		return a.generateErrorResponse(err, req)
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       userSession.Context,
	}
}

func (a *App) handlerUserLogout(w http.ResponseWriter, req *http.Request) rooter.Response {
	deleteLoginCookie(w)

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
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
