package api

import (
	"net/http"

	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
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

func (a *App) handlerTest(req *http.Request) rooter.Response {
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
