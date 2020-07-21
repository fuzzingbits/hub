package rooter

import "net/http"

// Route is an API Route
type Route struct {
	Path    string
	Handler http.Handler
	Params  []Param
}

// Param is a route paramerter
type Param struct {
	Name     string
	Required bool
}

// RegisterRoutes for a given mux
func RegisterRoutes(mux *http.ServeMux, routes []Route) {
	for _, route := range routes {
		mux.Handle(route.Path, route.Handler)
	}
}
