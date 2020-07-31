package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

func TestSuccessfulRoutes2(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	// Insert test data
	request := entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	}

	requestBytes, _ := json.Marshal(request)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                   "test setup success",
			Method:                 http.MethodPost,
			URL:                    "/api/server/setup",
			Body:                   bytes.NewReader(requestBytes),
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:                "test setup already setup",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                bytes.NewReader(requestBytes),
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: responseServerAlreadySetup.Bytes(),
		},
	})
}
func TestServerStatus(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                   "status status",
			Method:                 http.MethodGet,
			URL:                    "/api/server/status",
			Body:                   nil,
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
	})
	c.UserProviderError = errors.New("foobar")
	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                   "test failed status",
			Method:                 http.MethodGet,
			URL:                    "/api/server/status",
			Body:                   nil,
			TargetStatusCode:       http.StatusInternalServerError,
			SkipResponseBytesCheck: true,
		},
	})
}

func TestSuccessfulRoutes(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	// Insert test data
	targetSession, _ := s.SetupServer(
		entity.CreateUserRequest{
			FirstName: "Testy",
			LastName:  "McTestPants",
			Username:  "testy",
			Email:     "testy@example.com",
			Password:  "Password123",
		},
	)

	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "Password123",
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name: "test me logged in",
			URL:  "/api/user/me",
			RequestMod: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:    session.CookieName,
					Value:   targetSession.Token,
					Expires: time.Now().Add(time.Minute),
				})
			},
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       targetSession.Context,
			}.Bytes(),
		},
		{
			Name:                "test no cookie me request",
			TargetStatusCode:    http.StatusOK,
			URL:                 "/api/user/me",
			TargetResponseBytes: responseMissingValidSession.Bytes(),
		},
		{
			Name:             "test login success",
			Method:           http.MethodPost,
			URL:              "/api/user/login",
			Body:             bytes.NewReader(loginRequestBytes),
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       targetSession.Context,
			}.Bytes(),
		},
	})
}

func TestFailedRoutes(t *testing.T) {
	container := container.NewMockable()
	service := hub.NewService(&hubconfig.Config{RollbarToken: "foobar-fake-token"}, container)
	mux := http.NewServeMux()
	RegisterRoutes(mux, service)

	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "Password123",
	}

	container.SessionProviderError = errors.New("foobar")

	loginRequestBytes, _ := json.Marshal(loginRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                "test test route",
			TargetStatusCode:    http.StatusInternalServerError,
			URL:                 "/api/user/me",
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "test login error 1",
			Method:              http.MethodPost,
			URL:                 "/api/user/login",
			Body:                bytes.NewReader(loginRequestBytes),
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: responseInvalidLogin.Bytes(),
		},
		{
			Name:                "test login error 2",
			Method:              http.MethodPost,
			URL:                 "/api/user/login",
			Body:                nil,
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "test setup",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                nil,
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
	})

	container.UserProviderError = errors.New("foobar")
	container.SessionProviderError = errors.New("foobar")

	request := entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	}

	requestBytes, _ := json.Marshal(request)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                "test me no user provider",
			TargetStatusCode:    http.StatusInternalServerError,
			URL:                 "/api/user/me",
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "test setup",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                nil,
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "test setup",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                bytes.NewReader(requestBytes),
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "test login success",
			Method:              http.MethodPost,
			URL:                 "/api/user/login",
			Body:                bytes.NewReader(loginRequestBytes),
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
	})
}
