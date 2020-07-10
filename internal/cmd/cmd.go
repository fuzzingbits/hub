package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/fuzzingbits/forge-wip/pkg/web"
	"github.com/fuzzingbits/hub/internal/api"
	"github.com/fuzzingbits/hub/internal/hub"
	"github.com/gobuffalo/packr"
)

// Run the hub command
func Run() {
	service, err := hub.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	server := getServer(service)
	log.Printf("Listening on: http://%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func getServer(service *hub.Service) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", getRootHandler(service))
	api.RegisterRoutes(mux, service)

	return &http.Server{
		Addr:         service.Config.Listen,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func getRootHandler(service *hub.Service) http.Handler {
	if service.Config.Dev {
		uiURL, _ := url.Parse(service.Config.DevUIProxyAddr)
		return httputil.NewSingleHostReverseProxy(uiURL)
	}

	uiFileSystem := packr.NewBox("../../dist")

	spaHandler := &web.SinglePageAppHandler{
		FileSystem:       uiFileSystem,
		FileName:         "index.html",
		DisableCSPHeader: true,
	}

	return &web.Handler{
		FileSystem:      uiFileSystem,
		NotFoundHandler: spaHandler,
		RootHandler:     spaHandler,
	}
}
