package web

import (
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
		ModResponse: func(w http.ResponseWriter, r *http.Request) {
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
			Name:                "test_spa_response_on_404",
			URL:                 "/fakepage",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: spaResponse,
		},
		{
			Name:                "test_root",
			URL:                 "/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: spaResponse,
			CustomResponseChecker: func(t *testing.T, r *http.Response) {
				targetCSP := "script-src 'self' 'sha256-QXZRmRPAsseuAgOGnvjVUJOnlHEzu25Ou1XhFOWnqyI='"
				actualCSP := r.Header.Get("Content-Security-Policy")
				if actualCSP != targetCSP {
					t.Fatalf(
						"%s returned: %s expected: %s",
						r.Request.URL.Path,
						actualCSP,
						targetCSP,
					)
				}

				targetTestHeader := "testValue"
				actualTestHeader := r.Header.Get("TestHeader")
				if actualTestHeader != targetTestHeader {
					t.Fatalf(
						"%s returned: %s expected: %s",
						r.Request.URL.Path,
						actualTestHeader,
						targetTestHeader,
					)
				}
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
