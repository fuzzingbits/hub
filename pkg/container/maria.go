package container

import (
	"github.com/jinzhu/gorm"
)

func (c *Production) getMariaClient() (*gorm.DB, error) {
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
	client, err := gorm.Open("mysql", c.config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	// TODO: configure the client here
	client.LogMode(false)

	// Save the successful connection
	c.mariaClient = client

	return c.mariaClient, nil
}
