package api

import (
	"net/http"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/util/forge/rooter"
	"github.com/fuzzingbits/hub/pkg/util/forge/rootertest"
)

func TestMiddlewareRecover(t *testing.T) {
	c := container.NewMockable()
	s := hub.NewService(&hubconfig.Config{}, c)
	a := &App{Service: s}

	mux := http.NewServeMux()

	mux.Handle(
		"/",
		a.middlewareRecovery(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					var s *hub.Service
					s.HTTPLogger.Print("this won't work")
				},
			),
		),
	)

	rootertest.Test(t, mux, []rootertest.TestCase{
		{
			Name:                "panic",
			URL:                 "/",
			TargetStatusCode:    http.StatusInternalServerError,
			TargetResponseBytes: rooter.ResponseInternalServerError.Bytes(),
		},
	})
}
