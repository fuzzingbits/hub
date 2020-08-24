package session

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/gomodule/redigo/redis"
)

// RedisProvider is a Redis SessionProvider
type RedisProvider struct {
	Connection redis.Conn
}

// Get a session by token
func (p *RedisProvider) Get(token string) (string, error) {
	result, err := p.Connection.Do("GET", token)
	if err != nil {
		return "", err
	}

	resultBytes, ok := result.([]byte)
	if !ok {
		return "", entity.ErrRecordNotFound
	}

	return string(resultBytes), nil
}

// Set a session by token
func (p *RedisProvider) Set(token string, userUUID string) error {
	if _, err := p.Connection.Do("SETEX", token, Duration.Seconds(), userUUID); err != nil {
		return err
	}

	return nil
}

// AutoMigrate the data connection
func (p *RedisProvider) AutoMigrate(clearExitstingData bool) error {
	if clearExitstingData {
		if _, err := p.Connection.Do("FLUSHALL"); err != nil {
			return err
		}
	}

	return nil
}
