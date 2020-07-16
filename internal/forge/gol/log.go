package gol

import (
	"log"
	"net/http"
)

// LogLogger is a logger for the standard library logger
type LogLogger struct {
	Logger *log.Logger
}

// Log a message just to the previously provided logger
func (r *LogLogger) Log(message string, v ...interface{}) {
	if r.Logger == nil {
		return
	}

	r.Logger.Printf(message, v...)
}

// Error reports an error
func (r *LogLogger) Error(err error) {
	r.Log(err.Error())
}

// ErrorRequest reports an error and a http request
func (r *LogLogger) ErrorRequest(err error, req *http.Request) {
	r.Log(err.Error())
}
