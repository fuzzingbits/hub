package container

import (
	"sync"

	"github.com/fuzzingbits/hub/internal/hubconfig"
	"github.com/fuzzingbits/hub/internal/provider/user"
	"github.com/fuzzingbits/hub/internal/provider/usersettings"
	redis "github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
)

// Container exposes data providers
type Container interface {
	// AutoMigrate the data connections
	AutoMigrate(devMode bool) error
	// UserProvider safety builds and returns the Provider
	UserProvider() (user.Provider, error)
	// UserSettingsProvider safety builds and returns the Provider
	UserSettingsProvider() (usersettings.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *hubconfig.Config
	// Providers
	userProvider         *user.DatabaseProvider
	userSettingsProvider *usersettings.DatabaseProvider
	// Clients
	mariaClient *gorm.DB
	mongoClient *mongo.Client
	redisClient *redis.Client
	// Mutex Locks
	userProviderMutex         *sync.Mutex
	userSettingsProviderMutex *sync.Mutex
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
		mariaClientMutex:          &sync.Mutex{},
		mongoClientMutex:          &sync.Mutex{},
		redisClientMutex:          &sync.Mutex{},
	}
}

// AutoMigrate the data connections
func (c *Production) AutoMigrate(devMode bool) error {

	userProvider, err := c.UserProvider()
	if err != nil {
		return err
	}

	userSettingProvider, err := c.UserSettingsProvider()
	if err != nil {
		return err
	}

	if err := c.userProvider.AutoMigrate(devMode); err != nil {
		return err
	}

	if devMode {
		c.createFixtures(
			userProvider,
			userSettingProvider,
		)
	}

	return nil
}
