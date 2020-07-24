package hub

import (
	"log"
	"os"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
)

// Service is the internal API of Hub
type Service struct {
	config      *hubconfig.Config
	container   container.Container
	Logger      *log.Logger
	HTTPLogger  *log.Logger
	AuditLogger *log.Logger
	DebugLogger *log.Logger
}

// NewService returns a production instance of the service
func NewService(newConfig *hubconfig.Config, newContainer container.Container) *Service {
	return &Service{
		config:      newConfig,
		container:   newContainer,
		Logger:      log.New(os.Stderr, "[HUB_STD] ", log.LstdFlags),
		HTTPLogger:  log.New(os.Stderr, "[HUB_HTTP] ", log.LstdFlags),
		AuditLogger: log.New(os.Stderr, "[HUB_AUDIT] ", log.LstdFlags),
		DebugLogger: log.New(os.Stderr, "[HUB_DEBUG] ", log.LstdFlags),
	}
}
