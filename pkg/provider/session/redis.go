package session

import (
	"context"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/go-redis/redis/v8"
)

// RedisProvider is a Redis SessionProvider
type RedisProvider struct {
	Client *redis.Client
}

// Get a session by token
func (p *RedisProvider) Get(token string) (entity.UserSession, error) {
	var session entity.UserSession

	result := p.Client.Get(context.TODO(), token)
	if err := result.Scan(&session); err != nil {
		return entity.UserSession{}, err
	}

	return session, nil
}

// Set a session by token
func (p *RedisProvider) Set(token string, session entity.UserSession) error {
	p.Client.Set(
		context.TODO(),
		token,
		session,
		Duration,
	)

	return nil
}
