package api

import (
	"net/http"

	"github.com/fuzzingbits/forge-wip/pkg/rooter"
	"github.com/fuzzingbits/hub/internal/entity"
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
			Path:    "/api/users/me",
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
		Data: entity.UserSession{
			User: entity.User{
				FirstName: "Aaron",
				LastName:  "Ellington",
			},
		},
	}
}
