package container

import "github.com/gomodule/redigo/redis"

func (c *Production) getRedisClient() (redis.Conn, error) {
	// If it's ready, just return it
	if c.redisClient != nil {
		return c.redisClient, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.redisClientMutex.Lock()
	defer c.redisClientMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.redisClient != nil {
		return c.redisClient, nil
	}

	// Try to make a connection
	client, err := redis.Dial(
		"tcp",
		c.config.CacheAddress,
		redis.DialPassword(c.config.CachePassword),
	)

	if err != nil {
		return nil, err
	}

	// TODO: configure the client here

	// Save the successful connection
	c.redisClient = client

	return c.redisClient, nil
}
