package gol

import "net/http"

// Logger interface for gol loggers
type Logger interface {
	Log(message string, v ...interface{})
	Error(err error)
	ErrorRequest(err error, req *http.Request)
}
