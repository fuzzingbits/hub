package container

import "database/sql"

func (c *Production) getMariaClient() (*sql.DB, error) {
	// If it's ready, just return it
	if c.mariaClient != nil {
		return c.mariaClient, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.mariaClientMutex.Lock()
	defer c.mariaClientMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.mariaClient != nil {
		return c.mariaClient, nil
	}

	// Try to make a connection
	client, err := sql.Open("mysql", c.config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	// TODO: configure the client here

	// Save the successful connection
	c.mariaClient = client

	return c.mariaClient, nil
}
