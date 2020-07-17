package hub

import (
	"net/http"

	"github.com/fuzzingbits/hub/internal/entity"
)

// CreateUser creates a User
func (s *Service) CreateUser(uuid string, firstName string, lastName string) (entity.UserSession, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	user, err := userProvider.Create(entity.User{
		UUID:      uuid,
		FirstName: firstName,
		LastName:  lastName,
	})
	if err != nil {
		return entity.UserSession{}, err
	}

	// Setup new UserSettings with defaults
	userSettings := entity.UserSettings{
		ThemeColor: "tomato",
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
