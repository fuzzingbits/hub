package container

import (
	"fmt"

	"github.com/fuzzingbits/hub/internal/forge/mockableprovider"
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
			Provider: mockableprovider.NewProvider(),
		},
		UserSettingsProviderValue: &usersettings.Mockable{
			Provider: mockableprovider.NewProvider(),
		},
	}
}

// AutoMigrate the data connection
func (m *Mockable) AutoMigrate(clearExitstingData bool) error {
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
