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
	s := hub.NewService(&hubconfig.Config{}, c)
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
			Name:                "no_request_body",
			Method:              http.MethodPost,
			URL:                 "/api/server/setup",
			Body:                nil,
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError().Bytes(),
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
