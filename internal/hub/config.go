package hub

import (
	"fmt"

	"github.com/fuzzingbits/forge-wip/pkg/config"
)

// Config for the HUB command line tool
type Config struct {
	Listen         string `env:"LISTEN"`
	DevUIProxyAddr string `end:"DEV_UI_PROXY_ADDR"`
	Dev            bool   `env:"DEV"`
}

func getConfig() (*Config, error) {
	configParser := config.Config{
		Providers: []config.Provider{
			config.ProviderEnvironment{},
		},
	}

	// Defaults are here
	c := &Config{
		Listen:         "0.0.0.0:2020",
		DevUIProxyAddr: "http://0.0.0.0:3000",
	}

	if err := configParser.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("Error parsing config: %w", err)
	}

	return c, nil
}
