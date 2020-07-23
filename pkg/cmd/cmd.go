package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/fuzzingbits/hub/pkg/api"
	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/util/forge/web"
	"github.com/gobuffalo/packr"
)

// App contains the required setup before running the app
type App struct {
	Config    *hubconfig.Config
	Container container.Container
	Service   *hub.Service
	Server    *http.Server
}

// Run the hub command
func Run() {
	app := App{}

	var err error
	if app.Config, err = hubconfig.GetConfig(); err != nil {
		log.Fatal(err)
	}

	app.Container = container.NewProduction(app.Config)
	app.Service = hub.NewService(app.Config, app.Container)
	server := getServer(app)
	go app.autoMigrate()

	log.Printf("Listening on: http://%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func getServer(app App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", getRootHandler(app))
	api.RegisterRoutes(mux, app.Service)

	return &http.Server{
		Addr:         app.Config.Listen,
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func getRootHandler(app App) http.Handler {
	if app.Config.DevProxyToNuxt {
		uiURL, _ := url.Parse(app.Config.DevUIProxyAddr)
		return httputil.NewSingleHostReverseProxy(uiURL)
	}

	uiFileSystem := packr.NewBox("../../dist")

	spaHandler := &web.SinglePageAppHandler{
		FileSystem: uiFileSystem,
		FileName:   "index.html",
		BaseCSPEntries: web.CSPEntries{
			Script: []string{"'self'"},
			Style:  []string{"'self'"},
		},
		ModResponse: func(w http.ResponseWriter) {
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("X-Content-Type-Options", "nosniff")
		},
	}

	return &web.Handler{
		FileSystem:      uiFileSystem,
		NotFoundHandler: spaHandler,
		RootHandler:     spaHandler,
	}
}

func (app App) autoMigrate() {
	var lastError error
	maxTryCount := 5
	postFailureWait := time.Second * 30

	for tryCount := 1; tryCount <= maxTryCount; tryCount++ {
		if lastError = app.Container.AutoMigrate(app.Config.DevClearExitstingData); lastError != nil {
			log.Printf("AutoMigrate Attempt Failed %d/%d: Waiting %.0f seconds before trying again...", tryCount, maxTryCount, postFailureWait.Seconds())
			time.Sleep(postFailureWait)
			continue
		}

		log.Printf("AutoMigrate Attempt Successful %d/%d", tryCount, maxTryCount)

		return
	}

	log.Printf("AutoMirgate Error: %s", lastError.Error())
}
