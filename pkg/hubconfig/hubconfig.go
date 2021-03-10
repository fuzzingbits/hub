package hubconfig

import (
	"fmt"

	"github.com/fuzzingbits/hub/pkg/util/forge/config"
)

// Config for the Hub command line tool
type Config struct {
	Listen           string `env:"LISTEN"`
	RollbarToken     string `env:"ROLLBAR_TOKEN"`
	DatabaseDSN      string `env:"DATABASE_DSN"`
	DocumentStoreDSN string `env:"DOCUMENT_STORE_DSN"`
	CacheAddress     string `env:"CACHE_ADDRESS"`
	CacheUsername    string `env:"CACHE_USERNAME"`
	CachePassword    string `env:"CACHE_PASSWORD"`
	// Development Parameters
	DevProxyToNuxt        bool   `env:"DEV_PROXY_TO_NUXT"`
	DevClearExitstingData bool   `env:"DEV_CLEAR_EXISTING_DATA"`
	DevUIProxyAddr        string `env:"DEV_UI_PROXY_ADDR"`
}

// GetConfig gets the config from the environment
func GetConfig() (*Config, error) {
	configParser := config.Config{
		Providers: []config.Provider{
			config.ProviderEnvironment{},
		},
	}

	// Default values are here
	c := &Config{
		Listen:           "0.0.0.0:8000",
		DevUIProxyAddr:   "http://0.0.0.0:3000",
		DatabaseDSN:      "root:justTheDevPassword@(127.0.0.1:2021)/hub?charset=utf8&parseTime=True&loc=Local",
		DocumentStoreDSN: "mongodb://root:justTheDevPassword@127.0.0.1:2022",
		CacheAddress:     "127.0.0.1:2023",
	}

	if err := configParser.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("Error parsing config: %w", err)
	}

	return c, nil
}
