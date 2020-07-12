package api

import (
	"net/http"

	"github.com/fuzzingbits/hub/internal/hub"
	"github.com/fuzzingbits/hub/internal/util/rooter"
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
			Path:    "/api/test",
			Handler: rooter.ResponseFunc(a.testHandler),
		},
	}

	rooter.RegisterRoutes(mux, routes)
}

func (a *App) testHandler(req *http.Request) rooter.Response {
	return rooter.Response{
		StatusCode: http.StatusOK,
		State:      true,
		Message:    "Hello!",
	}
}