package api

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/forge/rooter"
	"github.com/fuzzingbits/hub/internal/forge/rootertest"
	"github.com/fuzzingbits/hub/internal/hub"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestSuccessfulRoutes(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	// Insert test data
	targetSession, _ := s.CreateUser(
		entity.CreateUserRequest{
			FirstName: "Testy",
			LastName:  "McTestPants",
			Username:  "testy",
			Email:     "testy@example.com",
			Password:  "Password123",
		},
	)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name: "test test route",
			URL:  "/api/user/me",
			RequestMod: func(req *http.Request) {
				req.Header.Add("UUID", targetSession.User.UUID)
			},
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       targetSession,
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
			Name:             "test test route",
			TargetStatusCode: http.StatusOK,
			URL:              "/api/user/me",
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       nil,
			}.Bytes(),
		},
	})
}
