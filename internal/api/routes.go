package api

import (
	"net/http"

	"github.com/fuzzingbits/hub/internal/forge/rooter"
	"github.com/fuzzingbits/hub/internal/hub"
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
			Path:    "/api/server/status",
			Handler: rooter.ResponseFunc(a.testHandler),
		},
		{
			Path:    "/api/server/setup",
			Handler: rooter.ResponseFunc(a.testHandler),
		},
		{
			Path:    "/api/user/me",
			Handler: rooter.ResponseFunc(a.testHandler),
		},
	}

	rooter.RegisterRoutes(mux, routes)
}

func (a *App) testHandler(req *http.Request) rooter.Response {
	session, err := a.Service.GetCurrentSession(req)
	if err != nil {
		return rooter.Response{
			StatusCode: http.StatusOK,
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
