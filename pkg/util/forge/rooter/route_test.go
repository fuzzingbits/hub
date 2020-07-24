package rooter

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

func TestTestRoute(t *testing.T) {
	pathBasicTest := "/hello/world"
	responseBasicTest := Response{
		StatusCode: http.StatusOK,
		State:      true,
		Message:    "Hello, world!",
	}

	routes := []Route{
		{
			Path: pathBasicTest,
			Handler: ResponseFunc(func(r *http.Request) Response {
				return responseBasicTest
			}),
			Middleware: []Middleware{
				MiddlewareLogger,
			},
		},
	}

	mux := http.NewServeMux()
	RegisterRoutes(mux, routes, []Middleware{
		MiddlewareLogger,
	})

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                "Basic Route Test",
			TargetStatusCode:    http.StatusOK,
			URL:                 pathBasicTest,
			TargetResponseBytes: responseBasicTest.Bytes(),
		},
	})
}
