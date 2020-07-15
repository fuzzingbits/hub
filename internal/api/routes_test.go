package api

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/forge-wip/pkg/rooter"
	"github.com/fuzzingbits/forge-wip/pkg/rootertest"
	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/hub"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestSuccessfulRoutes(t *testing.T) {
	container := container.NewMockable()
	service := hub.NewService(&hubconfig.Config{}, container)
	mux := http.NewServeMux()
	RegisterRoutes(mux, service)

	// Insert test data
	targetUser := entity.User{
		UUID:      "313efbe9-173b-4a1b-9b5b-7b69d95a66b9",
		FirstName: "Testy",
		LastName:  "McTestPants",
	}

	container.UserProviderValue.Create(targetUser)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:       "test test route",
			StatusCode: http.StatusOK,
			URL:        "/api/users/me",
			ResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data: entity.UserSession{
					User: targetUser,
				},
			}.Bytes(),
		},
	})
}

func TestFailedRoutes(t *testing.T) {
	container := container.NewMockable()
	service := hub.NewService(&hubconfig.Config{}, container)
	mux := http.NewServeMux()
	RegisterRoutes(mux, service)

	// Test with no UserProvider
	container.UserProviderValue = nil

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:       "test test route",
			StatusCode: http.StatusOK,
			URL:        "/api/users/me",
			ResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       nil,
			}.Bytes(),
		},
	})
}
