package web

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

func TestSinglePageAppHandler(t *testing.T) {
	spaResponse := []byte("Hello, world!<script>console.log('foobar');</script>\n")
	fs := http.Dir("handler_test_files")

	singlePageAppHandler := &SinglePageAppHandler{
		FileSystem: fs,
		FileName:   "spa.html",
		BaseCSPEntries: CSPEntries{
			Script: []string{"'self'"},
		},
		ModResponse: func(w http.ResponseWriter) {
			w.Header().Set("TestHeader", "testValue")
		},
	}

	testHandler := &Handler{
		RootHandler:     singlePageAppHandler,
		NotFoundHandler: singlePageAppHandler,
		FileSystem:      fs,
	}

	testCases := []rootertest.TestCase{
		{
			URL:                 "/fakepage",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: spaResponse,
		},
		{
			URL:                 "/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: spaResponse,
			ResponseChecker: func(r *http.Response) error {
				targetCSP := "script-src 'self' 'sha256-QXZRmRPAsseuAgOGnvjVUJOnlHEzu25Ou1XhFOWnqyI='"
				actualCSP := r.Header.Get("Content-Security-Policy")
				if actualCSP != targetCSP {
					return fmt.Errorf(
						"%s returned: %s expected: %s",
						r.Request.URL.Path,
						actualCSP,
						targetCSP,
					)
				}

				targetTestHeader := "testValue"
				actualTestHeader := r.Header.Get("TestHeader")
				if actualTestHeader != targetTestHeader {
					return fmt.Errorf(
						"%s returned: %s expected: %s",
						r.Request.URL.Path,
						actualTestHeader,
						targetTestHeader,
					)
				}
				return nil
			},
		},
	}

	rootertest.Test(t, testHandler, testCases)
}

func TestSinglePageAppHandlerIndexNotFound(t *testing.T) {
	fs := http.Dir("handler_test_files")

	singlePageAppHandler := &SinglePageAppHandler{
		FileSystem: fs,
		FileName:   "notfound.html",
	}

	testHandler := &Handler{
		RootHandler:     singlePageAppHandler,
		NotFoundHandler: singlePageAppHandler,
		FileSystem:      fs,
	}

	testCases := []rootertest.TestCase{
		{
			URL:                 "/fakepage",
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: []byte("open handler_test_files/notfound.html: no such file or directory\n"),
		},
	}

	rootertest.Test(t, testHandler, testCases)
}
