package hub

import "github.com/fuzzingbits/hub/internal/container"

// Service is the internal API of HUB
type Service struct {
	Config    *Config
	container *container.Container
}

// NewProduction returns a production instance of the service
func NewProduction() (*Service, error) {
	config, err := getConfig()
	if err != nil {
		return nil, err
	}

	return &Service{
		Config:    config,
		container: &container.Container{},
	}, nil
}
