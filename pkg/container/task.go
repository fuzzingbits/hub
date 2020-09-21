package container

import (
	"github.com/fuzzingbits/hub/pkg/provider/task"
)

// TaskProvider safely builds and returns the Provider
func (c *Production) TaskProvider() (task.Provider, error) {
	// If it's ready, just return it
	if c.taskProvider != nil {
		return c.taskProvider, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.taskProviderMutex.Lock()
	defer c.taskProviderMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.taskProvider != nil {
		return c.taskProvider, nil
	}

	// Try to make a connection
	db, err := c.getMariaClient()
	if err != nil {
		return nil, err
	}

	// Create and save the provder
	c.taskProvider = &task.DatabaseProvider{
		Database: db,
	}

	return c.taskProvider, nil
}
