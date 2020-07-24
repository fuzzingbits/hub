package rooter

import "net/http"

// Route is an API Route
type Route struct {
	Path       string
	Handler    http.Handler
	Params     []Param
	Middleware []Middleware
}

// Param is a route paramerter
type Param struct {
	Name     string
	Required bool
}

// RegisterRoutes for a given mux
func RegisterRoutes(mux *http.ServeMux, routes []Route, middlewares []Middleware) {
	for _, route := range routes {
		// Add route specific middlewares
		for _, middleware := range route.Middleware {
			route.Handler = middleware(route.Handler)
		}

		// Add the global middleares
		for _, middleware := range middlewares {
			route.Handler = middleware(route.Handler)
		}

		mux.Handle(route.Path, route.Handler)
	}
}
