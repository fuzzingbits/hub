package container

import (
	"github.com/fuzzingbits/hub/pkg/provider/usersettings"
)

// UserSettingsProvider safety builds and returns the Provider
func (c *Production) UserSettingsProvider() (usersettings.Provider, error) {
	// If it's ready, just return it
	if c.userSettingsProvider != nil {
		return c.userSettingsProvider, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.userProviderMutex.Lock()
	defer c.userProviderMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.userSettingsProvider != nil {
		return c.userSettingsProvider, nil
	}

	// Try to make a connection
	db, err := c.getMongoClient()
	if err != nil {
		return nil, err
	}

	// Create and save the provder
	c.userSettingsProvider = &usersettings.DatabaseProvider{
		Collection: db.Database("hub").Collection("user_settings"),
	}

	return c.userSettingsProvider, nil
}
