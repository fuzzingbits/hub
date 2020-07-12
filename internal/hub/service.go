package hub

import (
	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

// Service is the internal API of HUB
type Service struct {
	config    *hubconfig.Config
	container *container.Container
}

// NewService returns a production instance of the service
func NewService(newConfig *hubconfig.Config, newContainer *container.Container) *Service {
	return &Service{
		config:    newConfig,
		container: newContainer,
	}
}
