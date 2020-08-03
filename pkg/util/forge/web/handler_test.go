package web

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

func TestHandler(t *testing.T) {
	testHandler := &Handler{
		FileSystem:  http.Dir("handler_test_files"),
		ModResponse: func(w http.ResponseWriter, r *http.Request) {},
	}

	testCases := []rootertest.TestCase{
		{
			Name:                "test_root",
			URL:                 "/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("root index\n"),
		},
		{
			Name:                "test_static_file",
			URL:                 "/test.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("root test\n"),
		},
		{
			Name:                "test_404_dir",
			URL:                 "/testdir",
			TargetStatusCode:    http.StatusNotFound,
			TargetResponseBytes: []byte("404 page not found\n"),
		},
		{
			Name:                "test_404_dir_with_trailing_slash",
			URL:                 "/testdir/",
			TargetStatusCode:    http.StatusNotFound,
			TargetResponseBytes: []byte("404 page not found\n"),
		},
		{
			Name:                "test_multiple_sub_directory_index",
			URL:                 "/testdir/testdir2/index.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			Name:                "test_sub_directory_index_with_trailing_slash",
			URL:                 "/testdir/testdir2/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			Name:                "test_sub_directory_index",
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
			Name:                "test_root",
			URL:                 "/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: customHandlerResponse,
		},
		{
			Name:                "test_root_static_file",
			URL:                 "/test.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("root test\n"),
		},
		{
			Name:                "test_sub_directory",
			URL:                 "/testdir",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: customHandlerResponse,
		},
		{
			Name:                "test_sub_directory_with_trailing_slash",
			URL:                 "/testdir/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: customHandlerResponse,
		},
		{
			Name:                "test_multiple_sub_directory_index",
			URL:                 "/testdir/testdir2/index.html",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			Name:                "test_multiple_sub_directory_with_trailing_slash",
			URL:                 "/testdir/testdir2/",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
		{
			Name:                "test_multiple_sub_directory",
			URL:                 "/testdir/testdir2",
			TargetStatusCode:    http.StatusOK,
			TargetResponseBytes: []byte("testdir testdir2 index\n"),
		},
	}

	rootertest.Test(t, testHandler, testCases)
}
