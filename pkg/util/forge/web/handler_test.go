package web

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

func TestHandler(t *testing.T) {
	testHandler := &Handler{
		FileSystem: http.Dir("handler_test_files"),
	}

	testCases := []rootertest.TestCase{
		{
			URL:                 "/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("root index\n"),
		},
		{
			URL:                 "/test.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("root test\n"),
		},
		{
			URL:                 "/testdir",
			TargetStatusCode:    http.StatusNotFound,
			TargetResponseBytes: []byte("404 page not found\n"),
		},
		{
			URL:                 "/testdir/",
			TargetStatusCode:    http.StatusNotFound,
			TargetResponseBytes: []byte("404 page not found\n"),
		},
		{
			URL:                 "/testdir/testdir2/index.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			URL:                 "/testdir/testdir2/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			URL:                 "/testdir/testdir2",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
	}

	rootertest.Test(t, testHandler, testCases)
}

func TestHandlerCustomHandler(t *testing.T) {
	customHandlerResponse := []byte("custom root")
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write(customHandlerResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	testHandler := &Handler{
		FileSystem:      http.Dir("handler_test_files"),
		RootHandler:     http.HandlerFunc(rootHandler),
		NotFoundHandler: http.HandlerFunc(rootHandler),
	}

	testCases := []rootertest.TestCase{
		{
			URL:                 "/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: customHandlerResponse,
		},
		{
			URL:                 "/test.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("root test\n"),
		},
		{
			URL:                 "/testdir",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: customHandlerResponse,
		},
		{
			URL:                 "/testdir/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: customHandlerResponse,
		},
		{
			URL:                 "/testdir/testdir2/index.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			URL:                 "/testdir/testdir2/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			URL:                 "/testdir/testdir2",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
	}

	rootertest.Test(t, testHandler, testCases)
}
