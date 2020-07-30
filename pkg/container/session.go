package container

import (
	"github.com/fuzzingbits/hub/pkg/provider/session"
)

// SessionProvider safely builds and returns the Provider
func (c *Production) SessionProvider() (session.Provider, error) {
	// If it's ready, just return it
	if c.sessionProvider != nil {
		return c.sessionProvider, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.sessionProviderMutex.Lock()
	defer c.sessionProviderMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.sessionProvider != nil {
		return c.sessionProvider, nil
	}

	// Try to make a connection
	client, err := c.getRedisClient()
	if err != nil {
		return nil, err
	}

	// Create and save the provder
	c.sessionProvider = &session.RedisProvider{
		Client: client,
	}

	return c.sessionProvider, nil
}
