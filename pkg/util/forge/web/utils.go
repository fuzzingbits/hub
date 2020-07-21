package web

import (
	"net/http"
)

const defaultDirectoryIndex = "index.html"

func fileExists(fs http.FileSystem, path string) bool {
	file, err := fs.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	return true
}
