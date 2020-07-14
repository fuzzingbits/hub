package api

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/forge-wip/pkg/rooter"
	"github.com/fuzzingbits/forge-wip/pkg/rootertest"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/hub"
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
				Data: entity.UserSession{
					User: entity.User{
						FirstName: "Aaron",
						LastName:  "Ellington",
					},
				},
			}.Bytes(),
		},
	})
}
