package api

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/internal/hub"
	"github.com/fuzzingbits/hub/internal/util/rooter"
	"github.com/fuzzingbits/hub/internal/util/rootertest"
)

func TestTestRoute(t *testing.T) {
	mux := http.NewServeMux()
	service := &hub.Service{}
	RegisterRoutes(mux, service)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:       "test test route",
			StatusCode: http.StatusOK,
			URL:        "/api/test",
			ResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Message:    "Hello!",
			}.Bytes(),
		},
	})
}