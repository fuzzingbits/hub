package container

import (
	"sync"

	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/provider/user"
	"github.com/fuzzingbits/hub/pkg/provider/usersettings"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container exposes data providers
type Container interface {
	// AutoMigrate the data connections
	AutoMigrate(clearExitstingData bool) error
	// UserProvider safely builds and returns the Provider
	UserProvider() (user.Provider, error)
	// SessionProvider safely builds and returns the Provider
	SessionProvider() (session.Provider, error)
	// UserSettingsProvider safely builds and returns the Provider
	UserSettingsProvider() (usersettings.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *hubconfig.Config
	// Providers
	userProvider         *user.DatabaseProvider
	userSettingsProvider *usersettings.DatabaseProvider
	sessionProvider      *session.RedisProvider
	// Clients
	mariaClient *gorm.DB
	mongoClient *mongo.Client
	redisClient redis.Conn
	// Mutex Locks
	userProviderMutex         *sync.Mutex
	userSettingsProviderMutex *sync.Mutex
	sessionProviderMutex      *sync.Mutex
	mariaClientMutex          *sync.Mutex
	mongoClientMutex          *sync.Mutex
	redisClientMutex          *sync.Mutex
}

// NewProduction builds a container with all of the config
func NewProduction(hubConfig *hubconfig.Config) Container {
	return &Production{
		config:                    hubConfig,
		userProviderMutex:         &sync.Mutex{},
		userSettingsProviderMutex: &sync.Mutex{},
		sessionProviderMutex:      &sync.Mutex{},
		mariaClientMutex:          &sync.Mutex{},
		mongoClientMutex:          &sync.Mutex{},
		redisClientMutex:          &sync.Mutex{},
	}
}

// AutoMigrate the data connections
func (c *Production) AutoMigrate(clearExitstingData bool) error {
	if _, err := c.UserProvider(); err != nil {
		return err
	}

	if _, err := c.UserSettingsProvider(); err != nil {
		return err
	}

	if _, err := c.SessionProvider(); err != nil {
		return err
	}

	if err := autoMigrateAll([]dataProvider{
		c.userProvider,
		c.userSettingsProvider,
		c.sessionProvider,
	}, clearExitstingData); err != nil {
		return err
	}

	return nil
}

type dataProvider interface {
	AutoMigrate(clearExitstingData bool) error
}

func autoMigrateAll(providers []dataProvider, clearExitstingData bool) error {
	for _, provider := range providers {
		if err := provider.AutoMigrate(clearExitstingData); err != nil {
			return err
		}
	}

	return nil
}
