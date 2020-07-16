package rooter

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/internal/forge/rootertest"
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
		},
	}

	mux := http.NewServeMux()
	RegisterRoutes(mux, routes)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                "Basic Route Test",
			TargetStatusCode:    http.StatusOK,
			URL:                 pathBasicTest,
			TargetResponseBytes: responseBasicTest.Bytes(),
		},
	})
}
