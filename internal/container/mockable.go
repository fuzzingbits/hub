package container

import (
	"fmt"
	"sync"

	"github.com/fuzzingbits/hub/internal/provider/user"
	"github.com/fuzzingbits/hub/internal/provider/usersettings"
)

// Mockable container struct
type Mockable struct {
	UserProviderValue         *user.Mockable
	UserProviderError         error
	UserSettingsProviderValue *usersettings.Mockable
	UserSettingsProviderError error
}

// NewMockable builds a new mockable container
func NewMockable() *Mockable {
	return &Mockable{
		UserProviderValue: &user.Mockable{
			Mutex: &sync.Mutex{},
		},
		UserSettingsProviderValue: &usersettings.Mockable{
			Mutex: &sync.Mutex{},
		},
	}
}

// AutoMigrate the data connection
func (m *Mockable) AutoMigrate(devMode bool) error {
	return nil
}

// UserProvider safety builds and returns the Provider
func (m *Mockable) UserProvider() (user.Provider, error) {
	if m.UserProviderValue == nil {
		return nil, fmt.Errorf("error getting the mockable user.Provider")
	}

	return m.UserProviderValue, m.UserProviderError
}

// UserSettingsProvider safety builds and returns the Provider
func (m *Mockable) UserSettingsProvider() (usersettings.Provider, error) {
	if m.UserSettingsProviderValue == nil {
		return nil, fmt.Errorf("error getting the mockable usersettings.Provider")
	}

	return m.UserSettingsProviderValue, m.UserSettingsProviderError
}
