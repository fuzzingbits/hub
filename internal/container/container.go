package container

import (
	"database/sql"
	"sync"

	"github.com/fuzzingbits/hub/internal/hubconfig"
	redis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container for our external connections
type Container struct {
	config *hubconfig.Config
	// Clients
	mariaClient *sql.DB
	mongoClient *mongo.Client
	redisClient *redis.Client
	// Mutex Locks
	mariaClientMutex *sync.Mutex
	mongoClientMutex *sync.Mutex
	redisClientMutex *sync.Mutex
}

// NewProduction builds a container with all of the config
func NewProduction(hubConfig *hubconfig.Config) *Container {
	return &Container{
		config: hubConfig,
	}
}
