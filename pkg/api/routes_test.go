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
var testUserCreateRequest = entity.UserCreateRequest{
	FirstName: "Testy",
	LastName:  "McTestPants",
	Email:     "testy@example.com",
	Password:  "Password123",
}

var testUserLoginRequest = entity.UserLoginRequest{
	Email:    testUserCreateRequest.Email,
	Password: testUserCreateRequest.Password,
}

func TestServerSetup(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{RollbarToken: "FAKE_TOKEN"}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	testUserCreateRequestBytes, _ := json.Marshal(testUserCreateRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                   "success",
			Method:                 http.MethodPost,
			URL:                    RouteServerSetup,
			Body:                   bytes.NewReader(testUserCreateRequestBytes),
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:                "already_setup",
			Method:              http.MethodPost,
			URL:                 RouteServerSetup,
			Body:                bytes.NewReader(testUserCreateRequestBytes),
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: ResponseServerAlreadySetup.Bytes(),
		},
		{
			Name:                "no_body",
			Method:              http.MethodPost,
			URL:                 RouteServerSetup,
			Body:                nil,
			TargetStatusCode:    http.StatusBadRequest,
			TargetResponseBytes: rooter.ResponseBadRequest.Bytes(),
		},
		{
			Name:   "server_error",
			Method: http.MethodPost,
			URL:    RouteServerSetup,
			Body:   bytes.NewReader(testUserCreateRequestBytes),
			RequestMod: func(r *http.Request) {
				c.UserProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError.Bytes(),
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
			URL:                    RouteServerStatus,
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "server_error",
			Method: http.MethodGet,
			URL:    RouteServerStatus,
			RequestMod: func(r *http.Request) {
				c.UserProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError.Bytes(),
		},
	})
}

func TestFavicon(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                   "success",
			Method:                 http.MethodGet,
			URL:                    RouteFavicon,
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
	})
}

func TestUserMe(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testUserCreateRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:   "success",
			Method: http.MethodGet,
			URL:    RouteUserMe,
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
			URL:    RouteUserMe,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: "FAKE_COOKIE_VALUE",
				})
			},
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
			}.Bytes(),
		},
		{
			Name:   "server_error",
			Method: http.MethodGet,
			URL:    RouteUserMe,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
				c.SessionProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError.Bytes(),
		},
		{
			Name:             "missing_cookie",
			Method:           http.MethodGet,
			URL:              RouteUserMe,
			TargetStatusCode: http.StatusOK,
			TargetResponseBytes: rooter.Response{
				StatusCode: http.StatusOK,
				State:      true,
			}.Bytes(),
		},
	})
}

func TestUserNew(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testUserCreateRequest)
	testUserCreateRequestBytes, _ := json.Marshal(entity.UserCreateRequest{
		Email: "foobar@example.com",
	})

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:   "create_user_error",
			Method: http.MethodPost,
			URL:    RouteUserNew,
			Body:   bytes.NewReader(testUserCreateRequestBytes),
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
				c.UserProviderValue.Provider.CreateError = errors.New("foobar")
			},
			TargetStatusCode:       http.StatusInternalServerError,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "no_body",
			Method: http.MethodPost,
			URL:    RouteUserNew,
			Body:   nil,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:       http.StatusBadRequest,
			SkipResponseBytesCheck: true,
		},
		{
			Name:                   "no_auth",
			Method:                 http.MethodPost,
			URL:                    RouteUserNew,
			TargetStatusCode:       http.StatusUnauthorized,
			SkipResponseBytesCheck: true,
		},
	})
}

func TestUserLogin(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testUserCreateRequest)

	loginRequestBytes, _ := json.Marshal(testUserLoginRequest)
	loginBadRequestBytes, _ := json.Marshal(entity.UserLoginRequest{
		Email:    "bad_email",
		Password: "bad_password",
	})

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:             "success",
			Method:           http.MethodPost,
			URL:              RouteUserLogin,
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
			URL:                 RouteUserLogin,
			Body:                bytes.NewReader(loginBadRequestBytes),
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: ResponseInvalidLogin.Bytes(),
		},
		{
			Name:   "server_error",
			Method: http.MethodPost,
			URL:    RouteUserLogin,
			Body:   bytes.NewReader(loginRequestBytes),
			RequestMod: func(r *http.Request) {
				c.SessionProviderError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError.Bytes(),
		},
		{
			Name:                "no_body",
			Method:              http.MethodPost,
			URL:                 RouteUserLogin,
			TargetStatusCode:    http.StatusBadRequest,
			TargetResponseBytes: rooter.ResponseBadRequest.Bytes(),
		},
	})
}

func TestUserDelete(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testUserCreateRequest)
	userContext, _ := s.CreateUser(entity.UserCreateRequest{
		Email: "foobar@example.com",
	})
	payloadBytes, _ := json.Marshal(entity.UserDeleteRequest{
		UUID: userContext.User.UUID,
	})

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:   "success",
			Method: http.MethodPost,
			URL:    RouteUserDelete,
			Body:   bytes.NewReader(payloadBytes),
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "not_found",
			Method: http.MethodPost,
			URL:    RouteUserDelete,
			Body:   bytes.NewReader(payloadBytes),
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: ResponseRecordNotFound.Bytes(),
		},
		{
			Name:   "no_body",
			Method: http.MethodPost,
			URL:    RouteUserDelete,
			Body:   nil,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:    http.StatusBadRequest,
			TargetResponseBytes: rooter.ResponseBadRequest.Bytes(),
		},
		{
			Name:                "no_auth",
			Method:              http.MethodPost,
			URL:                 RouteUserDelete,
			Body:                nil,
			TargetStatusCode:    http.StatusUnauthorized,
			TargetResponseBytes: rooter.ResponseUnauthorized.Bytes(),
		},
	})
}

func TestUserUpdate(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testUserCreateRequest)
	payloadBytes, _ := json.Marshal(entity.UserUpdateRequest{
		UUID: userSession.Context.User.UUID,
	})

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:   "success",
			Method: http.MethodPost,
			URL:    RouteUserUpdate,
			Body:   bytes.NewReader(payloadBytes),
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "error_saving_session",
			Method: http.MethodPost,
			URL:    RouteUserUpdate,
			Body:   bytes.NewReader(payloadBytes),
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
				c.SessionProviderValue.Provider.CreateError = errors.New("foobar")
			},
			TargetStatusCode:       http.StatusInternalServerError,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "not_found",
			Method: http.MethodPost,
			URL:    RouteUserUpdate,
			Body:   bytes.NewReader(payloadBytes),
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
				c.UserProviderValue.Provider.UpdateError = errors.New("foobar")
			},
			TargetStatusCode:       http.StatusInternalServerError,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "no_body",
			Method: http.MethodPost,
			URL:    RouteUserUpdate,
			Body:   nil,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:    http.StatusBadRequest,
			TargetResponseBytes: rooter.ResponseBadRequest.Bytes(),
		},
		{
			Name:                "no_auth",
			Method:              http.MethodPost,
			URL:                 RouteUserUpdate,
			Body:                nil,
			TargetStatusCode:    http.StatusUnauthorized,
			TargetResponseBytes: rooter.ResponseUnauthorized.Bytes(),
		},
	})
}

func TestUserList(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	mux := http.NewServeMux()
	RegisterRoutes(mux, s)

	userSession, _ := s.SetupServer(testUserCreateRequest)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:   "success",
			Method: http.MethodGet,
			URL:    RouteUserList,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
			},
			TargetStatusCode:       http.StatusOK,
			SkipResponseBytesCheck: true,
		},
		{
			Name:   "service_error",
			Method: http.MethodGet,
			URL:    RouteUserList,
			Body:   nil,
			RequestMod: func(r *http.Request) {
				r.AddCookie(&http.Cookie{
					Name:  session.CookieName,
					Value: userSession.Token,
				})
				c.UserProviderValue.Provider.GetAllError = errors.New("foobar")
			},
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError.Bytes(),
		},
		{
			Name:                "no_auth",
			Method:              http.MethodGet,
			URL:                 RouteUserList,
			Body:                nil,
			TargetStatusCode:    http.StatusUnauthorized,
			TargetResponseBytes: rooter.ResponseUnauthorized.Bytes(),
		},
	})
}
