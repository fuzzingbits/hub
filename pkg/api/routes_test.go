package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

// Insert test data
var testCreateUserRequest = entity.CreateUserRequest{
	FirstName: "Testy",
	LastName:  "McTestPants",
	Username:  "testy",
	Email:     "testy@example.com",
	Password:  "Password123",
}

func TestServerSetup(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{RollbarToken: "FAKE_TOKEN"}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	testCreateUserRequestBytes, _ := json.Marshal(testCreateUserRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                   "success",
			Method:                 http.MethodPost,
			URL:                    "/api/server/setup",
			Body:                   bytes.NewReader(testCreateUserRequestBytes),
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:                "already_setup",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                bytes.NewReader(testCreateUserRequestBytes),
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: responseServerAlreadySetup.Bytes(),
		},
		{
			Name:                "no_body",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                nil,
			TargetStatusCode:    http.StatusBadRequest,
			TargetResponseBytes: rooter.ResponseBadRequest().Bytes(),
		},
		{
			Name:   "server_error",
			Method: http.MethodPost,
			URL:    "/api/server/setup",
			Body:   bytes.NewReader(testCreateUserRequestBytes),
			RequestMod: func(r *http.Request) {
				c.UserProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
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
			Name:                   "success",
			Method:                 http.MethodGet,
			URL:                    "/api/server/status",
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "server_error",
			Method: http.MethodGet,
			URL:    "/api/server/status",
			RequestMod: func(r *http.Request) {
				c.UserProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
	})
}

func TestUserMe(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testCreateUserRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:   "success",
			Method: http.MethodGet,
			URL:    "/api/user/me",
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       userSession.Context,
			}.Bytes(),
		},
		{
			Name:   "invalid_cookie",
			Method: http.MethodGet,
			URL:    "/api/user/me",
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: "FAKE_COOKIE_VALUE",
				})
			},
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: responseMissingValidSession.Bytes(),
		},
		{
			Name:   "server_error",
			Method: http.MethodGet,
			URL:    "/api/user/me",
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
				c.SessionProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "missing_cookie",
			Method:              http.MethodGet,
			URL:                 "/api/user/me",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: responseMissingValidSession.Bytes(),
		},
	})
}

func TestUserLogin(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testCreateUserRequest)

	loginRequest := entity.UserLoginRequest{
		Username: testCreateUserRequest.Username,
		Password: testCreateUserRequest.Password,
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)
	loginBadRequestBytes, _ := json.Marshal(entity.UserLoginRequest{
		Username: "bad_username",
		Password: "bad_password",
	})

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:             "success",
			Method:           http.MethodPost,
			URL:              "/api/user/login",
			Body:             bytes.NewReader(loginRequestBytes),
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
				Data:       userSession.Context,
			}.Bytes(),
		},
		{
			Name:                "bad_login",
			Method:              http.MethodPost,
			URL:                 "/api/user/login",
			Body:                bytes.NewReader(loginBadRequestBytes),
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: responseInvalidLogin.Bytes(),
		},
		{
			Name:   "server_error",
			Method: http.MethodPost,
			URL:    "/api/user/login",
			Body:   bytes.NewReader(loginRequestBytes),
			RequestMod: func(r *http.Request) {
				c.SessionProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
		},
		{
			Name:                "no_body",
			Method:              http.MethodPost,
			URL:                 "/api/user/login",
			TargetStatusCode:    http.StatusBadRequest,
			TargetResponseBytes: rooter.ResponseBadRequest().Bytes(),
		},
	})
}
