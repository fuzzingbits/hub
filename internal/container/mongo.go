package container

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (c *Container) getMongoClient() (*mongo.Client, error) {
	// If it's ready, just return it
	if c.mongoClient != nil {
		return c.mongoClient, nil
	}

	// Lock so nothing else can try to make a connection at the same time
	c.mongoClientMutex.Lock()
	defer c.mongoClientMutex.Unlock()

	// If someone else got it ready while we ere waiting, just use it
	if c.mongoClient != nil {
		return c.mongoClient, nil
	}

	// Try to make a connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.config.DocumentStoreDSN))
	if err != nil {
		return nil, err
	}

	// TODO: configure the client here

	// Save the successful connection
	c.mongoClient = client

	return c.mongoClient, nil
}
