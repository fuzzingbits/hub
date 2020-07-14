package container

import "github.com/go-redis/redis/v8"

func (c *Production) getRedisClient() (*redis.Client, error) {
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
	client := redis.NewClient(&redis.Options{
		Addr: c.config.CacheAddress,
	})

	// TODO: configure the client here

	// Save the successful connection
	c.redisClient = client

	return c.redisClient, nil
}
