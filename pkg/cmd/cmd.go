package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/fuzzingbits/hub/pkg/api"
	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/hub"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/util/forge/web"
	"github.com/gobuffalo/packr"
	"github.com/rollbar/rollbar-go"
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

	app.Service.Logger.Printf("Listening on: http://%s\n", server.Addr)
	app.logError(server.ListenAndServe())
	app.Service.Rollbar.Wait()
}

func getServer(app App) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", getRootHandler(app))
	api.RegisterRoutes(mux, app.Service)

	return &http.Server{
		Addr:         app.Config.Listen,
		Handler:      mux,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
}

func getRootHandler(app App) http.Handler {
	if app.Config.DevProxyToNuxt {
		uiURL, _ := url.Parse(app.Config.DevUIProxyAddr)
		return httputil.NewSingleHostReverseProxy(uiURL)
	}

	uiFileSystem := packr.NewBox("../../resources/dist")

	spaHandler := &web.SinglePageAppHandler{
		FileSystem: uiFileSystem,
		FileName:   "200.html",
		BaseCSPEntries: web.CSPEntries{
			Script: []string{"'self'"},
			Style:  []string{"'self'"},
		},
		ModResponse: func(w http.ResponseWriter, r *http.Request) {
			// Add security headers when serving the html app
			w.Header().Set("X-Frame-Options", "SAMEORIGIN")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// Don't cache the html page. It's the only thing that will change and it should always be small
			w.Header().Set("Cache-Control", "no-cache")
		},
	}

	return &web.Handler{
		FileSystem:      uiFileSystem,
		NotFoundHandler: spaHandler,
		RootHandler:     spaHandler,
		ModResponse: func(w http.ResponseWriter, r *http.Request) {
			// Make sure all nuxt content is cached because all the files names include a hash
			if strings.HasPrefix(r.URL.Path, "/_nuxt") {
				w.Header().Set("Cache-Control", "public, max-age=31536000")
			} else {
				w.Header().Set("Cache-Control", "public, max-age=3600")
			}
		},
	}
}

func (app App) autoMigrate() {
	var lastError error
	maxTryCount := 20
	postFailureWait := time.Second * 15

	for tryCount := 1; tryCount <= maxTryCount; tryCount++ {
		if lastError = app.Container.AutoMigrate(app.Config.DevClearExitstingData); lastError != nil {
			app.Service.Logger.Printf("AutoMigrate Attempt Failed %d/%d: Waiting %.0f seconds before trying again...", tryCount, maxTryCount, postFailureWait.Seconds())
			time.Sleep(postFailureWait)
			continue
		}

		app.Service.Logger.Printf("AutoMigrate Attempt Successful %d/%d", tryCount, maxTryCount)

		return
	}

	app.logError(fmt.Errorf("AutoMirgate Error: %s", lastError.Error()))
}

func (app App) logError(err error) {
	app.Service.ErrorLogger.Printf("AutoMirgate Error: %s", err.Error())
	if app.Service.Rollbar != nil {
		app.Service.Rollbar.Message(rollbar.ERR, err.Error())
	}
}
