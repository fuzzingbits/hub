package api

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
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
	service := hub.NewService(&hubconfig.Config{RollbarToken: "foobar-fake-token"}, container)
	mux := http.NewServeMux()
	RegisterRoutes(mux, service)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:             "test test route",
			TargetStatusCode: http.StatusOK,
			URL:              "/api/user/me",
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Message:    "you are not logged in",
				Data:       nil,
			}.Bytes(),
		},
	})

	container.UserProviderError = errors.New("foobar")

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                "test test route",
			TargetStatusCode:    http.StatusInternalServerError,
			URL:                 "/api/user/me",
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
	})
}
