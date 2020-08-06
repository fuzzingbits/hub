package rooter

import (
	"encoding/json"
	"net/http"
)

// Response is a standard response format to use for API responses
type Response struct {
	StatusCode int         `json:"-"`
	State      bool        `json:"state"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	ExtraData  interface{} `json:"extraData"`
}

// Bytes return
func (s Response) Bytes() []byte {
	b, _ := json.Marshal(s)
	return b
}

// ServeHTTP writes the response as JSON
func (s Response) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s.StatusCode)
	w.Write(s.Bytes())
}

// ResponseHandler is a simpler than http.HandlerFunc for enforcing the proper usage of Response
type ResponseHandler func(w http.ResponseWriter, r *http.Request) Response

// ResponseFunc is a shortcut for making sure a ResponseHandler can be used as a http.Handler
func ResponseFunc(s ResponseHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := s(w, r)
		response.ServeHTTP(w, r)
	})
}

// ResponseNotFound is a standard way of serving a 404
var ResponseNotFound = Response{
	StatusCode: http.StatusNotFound,
	Message:    "not found",
}

// ResponseUnauthorized is a standard way of serving a 401
var ResponseUnauthorized = Response{
	StatusCode: http.StatusUnauthorized,
	Message:    "unauthorized",
}

// ResponseBadRequest is a standard way of serving a 400
var ResponseBadRequest = Response{
	StatusCode: http.StatusBadRequest,
	Message:    "bad request",
}

// ResponseMethodNotAllowed is a standard way of serving a 405
var ResponseMethodNotAllowed = Response{
	StatusCode: http.StatusMethodNotAllowed,
	Message:    "method not allowed",
}

// ResponseInternalServerError is a standard way of serving a 500
var ResponseInternalServerError = Response{
	StatusCode: http.StatusInternalServerError,
	Message:    "internal server error",
}
