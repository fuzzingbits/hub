package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/fuzzingbits/forge-wip/pkg/web"
	"github.com/fuzzingbits/hub/internal/api"
	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/hub"
	"github.com/fuzzingbits/hub/internal/hubconfig"
	"github.com/gobuffalo/packr"
)

// App contains the required setup before running the app
type App struct {
	Config    *hubconfig.Config
	Container *container.Container
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
	app.Server = getServer(app)

	log.Printf("Listening on: http://%s\n", app.Server.Addr)
	log.Fatal(app.Server.ListenAndServe())
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
	if app.Config.Dev {
		uiURL, _ := url.Parse(app.Config.DevUIProxyAddr)
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