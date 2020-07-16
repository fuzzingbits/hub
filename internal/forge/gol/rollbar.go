package gol

import (
	"log"
	"net/http"

	"github.com/rollbar/rollbar-go"
)

// RollbarLogger is foobar
type RollbarLogger struct {
	Logger  *log.Logger
	Rollbar *rollbar.Client
}

// Log a message just to the previously provided logger
func (r *RollbarLogger) Log(message string, v ...interface{}) {
	if r.Logger == nil {
		return
	}

	r.Logger.Printf(message, v...)
}

// Error reports an error
func (r *RollbarLogger) Error(err error) {
	r.Log(err.Error())

	r.Rollbar.ErrorWithStackSkip(rollbar.ERR, err, 5)
	r.Rollbar.Wait()
}

// ErrorRequest reports an error and a http request
func (r *RollbarLogger) ErrorRequest(err error, req *http.Request) {
	r.Log(err.Error())

	r.Rollbar.RequestErrorWithStackSkip(rollbar.ERR, req, err, 5)
	r.Rollbar.Wait()
}
