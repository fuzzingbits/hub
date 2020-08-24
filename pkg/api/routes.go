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
	RouteFavicon      = "/favicon.svg"
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
			Payload:  entity.UserCreateRequest{},
			Response: entity.UserContext{},
		},
		{
			Path:     RouteUserNew,
			Handler:  rooter.ResponseFunc(a.handlerUserNew),
			Payload:  entity.UserCreateRequest{},
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
			Payload:  entity.UserUpdateRequest{},
			Response: entity.UserContext{},
		},
		{
			ExcludeFromTypeScript: true,
			Path:                  RouteFavicon,
			Handler:               http.HandlerFunc(a.handlerFavicon),
		},
	}
}

func (a *App) handlerFavicon(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("color")
	if color == "" {
		color = hub.DefaultThemeColorLight
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="512px" height="512px" viewBox="0 0 512 512" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
	<path d="M300.241565,11.7367846 L446.076299,95.934508 C473.181243,111.583554 489.878576,140.504184 489.878576,171.802277 L489.878576,340.197723 C489.878576,371.495816 473.181243,400.416446 446.076299,416.065492 L300.241565,500.263215 C273.136622,515.912262 239.741954,515.912262 212.637011,500.263215 L66.8022766,416.065492 C39.6973335,400.416446 23,371.495816 23,340.197723 L23,171.802277 C23,140.504184 39.6973335,111.583554 66.8022766,95.934508 L212.637011,11.7367846 C239.741954,-3.91226155 273.136622,-3.91226155 300.241565,11.7367846 Z" fill="` + color + `"></path>
</svg>`))
}

func (a *App) handlerUserUpdate(w http.ResponseWriter, req *http.Request) rooter.Response {
	// Require login
	userSession, err := a.authCheck(req)
	if err != nil {
		return a.generateErrorResponse(err, req)
	}

	// Get the payload
	var payload entity.UserUpdateRequest
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
	var payload entity.UserDeleteRequest
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
	var payload entity.UserCreateRequest
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
	var payload entity.UserCreateRequest
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
