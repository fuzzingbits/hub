package container

import (
	"github.com/fuzzingbits/hub/pkg/provider/user"
)

// UserProvider safety builds and returns the Provider
func (c *Production) UserProvider() (user.Provider, error) {
	// If it's ready, just return it
	if c.userProvider != nil {
		return c.userProvider, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.userProviderMutex.Lock()
	defer c.userProviderMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.userProvider != nil {
		return c.userProvider, nil
	}

	// Try to make a connection
	db, err := c.getMariaClient()
	if err != nil {
		return nil, err
	}

	// Create and save the provder
	c.userProvider = &user.DatabaseProvider{
		Database: db,
	}

	return c.userProvider, nil
}
