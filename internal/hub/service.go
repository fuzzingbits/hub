package hub

import (
	"net/http"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

// Service is the internal API of Hub
type Service struct {
	config    *hubconfig.Config
	container container.Container
}

// NewService returns a production instance of the service
func NewService(newConfig *hubconfig.Config, newContainer container.Container) *Service {
	return &Service{
		config:    newConfig,
		container: newContainer,
	}
}

// GetCurrentSession gets the current session
func (s *Service) GetCurrentSession(r *http.Request) (entity.UserSession, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	// TODO: create actual session management
	user, err := userProvider.GetUserByUUID("313efbe9-173b-4a1b-9a5b-7b69d95a66b9")
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettings, err := userSettingsProvider.GetByUUID(user.UUID)
	if err != nil {
		return entity.UserSession{}, err
	}

	return entity.UserSession{
		User:     user,
		Settings: userSettings,
	}, nil
}
