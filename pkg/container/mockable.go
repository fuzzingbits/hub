package container

import (
	"fmt"

	"github.com/fuzzingbits/hub/pkg/provider/routine"
	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/provider/user"
	"github.com/fuzzingbits/hub/pkg/provider/usersettings"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable container struct
type Mockable struct {
	UserProviderValue         *user.Mockable
	UserProviderError         error
	UserSettingsProviderValue *usersettings.Mockable
	UserSettingsProviderError error
	SessionProviderValue      *session.Mockable
	SessionProviderError      error
	RoutineProviderValue      *routine.Mockable
	RoutineProviderError      error
}

// NewMockable builds a new mockable container
func NewMockable() *Mockable {
	return &Mockable{
		UserProviderValue: &user.Mockable{
			Provider: mockableprovider.NewProvider(),
		},
		UserSettingsProviderValue: &usersettings.Mockable{
			Provider: mockableprovider.NewProvider(),
		},
		SessionProviderValue: &session.Mockable{
			Provider: mockableprovider.NewProvider(),
		},
		RoutineProviderValue: &routine.Mockable{
			Provider: mockableprovider.NewProvider(),
		},
	}
}

// AutoMigrate the data connection
func (m *Mockable) AutoMigrate(clearExitstingData bool) error {
	return nil
}

// UserProvider safely builds and returns the Provider
func (m *Mockable) UserProvider() (user.Provider, error) {
	if m.UserProviderValue == nil {
		return nil, fmt.Errorf("error getting the mockable user.Provider")
	}

	return m.UserProviderValue, m.UserProviderError
}

// UserSettingsProvider safely builds and returns the Provider
func (m *Mockable) UserSettingsProvider() (usersettings.Provider, error) {
	if m.UserSettingsProviderValue == nil {
		return nil, fmt.Errorf("error getting the mockable usersettings.Provider")
	}

	return m.UserSettingsProviderValue, m.UserSettingsProviderError
}

// SessionProvider safely builds and returns the Provider
func (m *Mockable) SessionProvider() (session.Provider, error) {
	if m.SessionProviderValue == nil {
		return nil, fmt.Errorf("error getting the mockable session.Provider")
	}

	return m.SessionProviderValue, m.SessionProviderError
}

// RoutineProvider safely builds and returns the Provider
func (m *Mockable) RoutineProvider() (routine.Provider, error) {
	if m.RoutineProviderValue == nil {
		return nil, fmt.Errorf("error getting the mockable session.Provider")
	}

	return m.RoutineProviderValue, m.RoutineProviderError
}
