package api

import (
	"errors"
	"net/http"

	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/provider/user"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/rollbar/rollbar-go"
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

	routes := []rooter.Route{
		{
			Path:    "/api/user/me",
			Handler: rooter.ResponseFunc(a.handlerTest),
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

func (a *App) handlerTest(req *http.Request) rooter.Response {
	session, err := a.Service.GetCurrentSession(req)
	if err != nil {
		if !errors.Is(err, user.ErrNotFound) {
			return a.serverError(err, req)
		}

		return rooter.Response{
			StatusCode: http.StatusOK,
			Message:    "you are not logged in",
			State:      true,
			Data:       nil,
		}
	}

	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Data:       session,
	}
}
