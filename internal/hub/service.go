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

// CreateUser creates a User
func (s *Service) CreateUser(firstName string, lastName string) (entity.UserSession, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	user, err := userProvider.Create(entity.User{
		FirstName: firstName,
		LastName:  lastName,
		UUID:      "something random",
	})
	if err != nil {
		return entity.UserSession{}, err
	}

	// Setup new UserSettings with defaults
	userSettings := entity.UserSettings{
		ThemeColor: "#00bfff",
	}

	if err := userSettingsProvider.Save(user.UUID, userSettings); err != nil {
		return entity.UserSession{}, err
	}

	return entity.UserSession{
		User:     user,
		Settings: userSettings,
	}, nil
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

	// TODO: create actual session management and stop doing this terrible thing
	user, err := userProvider.GetByUUID(r.Header.Get("UUID"))
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
