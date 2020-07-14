package container

import (
	"database/sql"
	"sync"

	"github.com/fuzzingbits/hub/internal/hubconfig"
	"github.com/fuzzingbits/hub/internal/provider/user"
	redis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container exposes data providers
type Container interface {
	UserProvider() (user.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *hubconfig.Config
	// Providers
	userProvider user.Provider
	// Clients
	mariaClient *sql.DB
	mongoClient *mongo.Client
	redisClient *redis.Client
	// Mutex Locks
	userProviderMutex *sync.Mutex
	mariaClientMutex  *sync.Mutex
	mongoClientMutex  *sync.Mutex
	redisClientMutex  *sync.Mutex
}

// NewProduction builds a container with all of the config
func NewProduction(hubConfig *hubconfig.Config) Container {
	return &Production{
		config: hubConfig,
	}
}
