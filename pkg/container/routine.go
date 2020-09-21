package container

import (
	"github.com/fuzzingbits/hub/pkg/provider/routine"
)

// RoutineProvider safely builds and returns the Provider
func (c *Production) RoutineProvider() (routine.Provider, error) {
	// If it's ready, just return it
	if c.routineProvider != nil {
		return c.routineProvider, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.routineProviderMutex.Lock()
	defer c.routineProviderMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.routineProvider != nil {
		return c.routineProvider, nil
	}

	// Try to make a connection
	db, err := c.getMariaClient()
	if err != nil {
		return nil, err
	}

	// Create and save the provder
	c.routineProvider = &routine.DatabaseProvider{
		Database: db,
	}

	return c.routineProvider, nil
}
