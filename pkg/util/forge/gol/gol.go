package gol

import "net/http"

// Logger interface for gol loggers
type Logger interface {
	Log(message string, v ...interface{})
	Error(err error)
	ErrorRequest(err error, req *http.Request)
}

// PanicOnError should only be used when an error should not ever happen but still should be accounted for
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
