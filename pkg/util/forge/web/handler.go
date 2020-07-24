package web

import (
	"net/http"
	"strings"
)

// Handler is the Forge Handler
type Handler struct {
	FileSystem      http.FileSystem
	RootHandler     http.Handler
	NotFoundHandler http.Handler
	ModResponse     func(http.ResponseWriter, *http.Request)
	fileServer      http.Handler
}

// ServeHTTP satisfies the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set the fileServer if it's not already set
	if h.fileServer == nil {
		h.fileServer = http.FileServer(h.FileSystem)
	}

	// If accessing root and a RootHandler is set, use it
	if r.URL.Path == "/" && h.RootHandler != nil {
		h.RootHandler.ServeHTTP(w, r)
		return
	}

	// Check state of the request and save some checks
	requestedFilename := r.URL.Path
	requestingDirectory := strings.HasSuffix(requestedFilename, "/")
	if requestingDirectory {
		requestedFilename += defaultDirectoryIndex
	}

	if h.ModResponse != nil {
		h.ModResponse(w, r)
	}

	// 404 Not Found Handling
	if !fileExists(h.FileSystem, requestedFilename) {
		h.notFound(w, r)
		return
	}

	h.fileServer.ServeHTTP(w, r)
}

func (h *Handler) notFound(w http.ResponseWriter, r *http.Request) bool {
	// Use custom 404 Not Found Handler if there is one
	if h.NotFoundHandler != nil {
		h.NotFoundHandler.ServeHTTP(w, r)
		return true
	}

	// Default to default 404 Not Found Handler
	http.NotFound(w, r)
	return true
}
