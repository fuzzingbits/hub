package api

import (
	"errors"
	"net/http"

	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/rollbar/rollbar-go"
)

// Errors
var (
	ErrUnauthorized = errors.New("unauthorized")
)

// Responses
var (
	ResponseServerAlreadySetup = rooter.Response{
		StatusCode: http.StatusOK,
		Message:    "Server Is Already Setup",
		State:      false,
	}
	ResponseMissingValidSession = rooter.Response{
		StatusCode: http.StatusOK,
		Message:    "Missing Valid Session",
		State:      false,
	}
	ResponseInvalidLogin = rooter.Response{
		StatusCode: http.StatusOK,
		State:      false,
		Message:    "Invlaid Login",
	}
)

func (a *App) generateErrorResponse(err error, r *http.Request) rooter.Response {
	// Return early if it's a known error
	if errors.Is(err, ErrUnauthorized) {
		return rooter.ResponseUnauthorized()
	}

	// Return early if it's a known error
	if errors.Is(err, hub.ErrInvalidLogin) {
		return ResponseInvalidLogin
	}

	// Return early if it's a known error
	if errors.Is(err, hub.ErrServerAlreadySetup) {
		return ResponseServerAlreadySetup
	}

	// Send to rollbar if Rollbar is configured
	if a.Service.Rollbar != nil {
		a.Service.Rollbar.RequestError(rollbar.ERR, r, err)
	}

	// Print the error to the ErrorLogger
	a.Service.ErrorLogger.Printf(
		"Request Error: %s %s - Err: %s",
		r.Method,
		r.URL.Path,
		err.Error(),
	)

	return rooter.ResponseInternalServerError()
}
