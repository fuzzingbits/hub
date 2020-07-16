package web

import (
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
)

// SinglePageAppHandler serves a file from the http.FileSystem and adds headers relevant to serving a secure
type SinglePageAppHandler struct {
	FileSystem       http.FileSystem
	FileName         string
	DisableCSPHeader bool
	BaseCSPEntries   CSPEntries
}

// ServeHTTP satisfies the http.Handler interface
func (s *SinglePageAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Open the file
	indexFile, err := s.FileSystem.Open(s.FileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer indexFile.Close()

	// Read the file contents
	fileContents, _ := ioutil.ReadAll(indexFile)

	// Add the appropriate http headers
	s.addHeaders(w, r, fileContents)

	// Write the file contents
	_, _ = w.Write(fileContents)
}

func (s *SinglePageAppHandler) addHeaders(w http.ResponseWriter, r *http.Request, fileContents []byte) {
	extension := filepath.Ext(s.FileName)
	w.Header().Set("Content-Type", mime.TypeByExtension(extension))

	w.Header().Set("Cache-Control", "public, max-age=31536000")

	if !s.DisableCSPHeader {
		csp := GenerateContentSecurityPolicy(fileContents, s.BaseCSPEntries)

		w.Header().Set("Content-Security-Policy", csp)
	}
}
