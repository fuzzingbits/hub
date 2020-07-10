package cmd

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/fuzzingbits/forge-wip/pkg/config"
	"github.com/fuzzingbits/forge-wip/pkg/web"
	"github.com/gobuffalo/packr"
)

// Config for the HUB command line tool
type Config struct {
	Listen string `env:"LISTEN"`
	Dev    bool   `env:"DEV"`
}

// Run the hub command
func Run() {
	c := getConfig()
	server := getServer(c)
	log.Printf("Listening on: http://%s\n", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func getConfig() Config {
	configParser := config.Config{
		Providers: []config.Provider{
			config.ProviderEnvironment{},
		},
	}

	// Defaults are here
	c := Config{
		Listen: "0.0.0.0:2020",
	}

	if err := configParser.Unmarshal(&c); err != nil {
		log.Fatal(err)
	}

	return c
}

func getServer(c Config) *http.Server {
	return &http.Server{
		Addr:         c.Listen,
		Handler:      getMux(c),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func getMux(c Config) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", getRootHandler(c))
	mux.HandleFunc("/api/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	return mux
}

func getRootHandler(c Config) http.Handler {
	if c.Dev {
		uiURL, _ := url.Parse("http://0.0.0.0:3000")
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
