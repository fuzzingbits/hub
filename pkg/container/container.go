package container

import (
	"sync"

	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/provider/routine"
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
	// RoutineProvider safely builds and returns the Provider
	RoutineProvider() (routine.Provider, error)
}

// Production is our production container for our external connections
type Production struct {
	config *hubconfig.Config
	// Maria Client
	mariaClient      *gorm.DB
	mariaClientMutex *sync.Mutex
	// Mongo Client
	mongoClient      *mongo.Client
	mongoClientMutex *sync.Mutex
	// Redis Client
	redisClient      redis.Conn
	redisClientMutex *sync.Mutex
	// User Provider
	userProvider      *user.DatabaseProvider
	userProviderMutex *sync.Mutex
	// User Settings Provider
	userSettingsProvider      *usersettings.DatabaseProvider
	userSettingsProviderMutex *sync.Mutex
	// Session Provider
	sessionProvider      *session.RedisProvider
	sessionProviderMutex *sync.Mutex
	// Routine Provider
	routineProvider      *routine.DatabaseProvider
	routineProviderMutex *sync.Mutex
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
		routineProviderMutex:      &sync.Mutex{},
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

	if _, err := c.RoutineProvider(); err != nil {
		return err
	}

	if err := autoMigrateAll([]dataProvider{
		c.userProvider,
		c.userSettingsProvider,
		c.sessionProvider,
		c.routineProvider,
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
