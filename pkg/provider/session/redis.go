package session

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"errors"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/gomodule/redigo/redis"
)

// RedisProvider is a Redis SessionProvider
type RedisProvider struct {
	Connection redis.Conn
}

// Get a session by token
func (p *RedisProvider) Get(token string) (entity.Session, error) {
	var session entity.Session

	result, err := p.Connection.Do("GET", token)
	if err != nil {
		return entity.Session{}, err
	}

	resultBytes, ok := result.([]byte)
	if !ok {
		return entity.Session{}, errors.New("no session found")
	}

	sessionBytes, err := hex.DecodeString(string(resultBytes))
	if err != nil {
		return entity.Session{}, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(sessionBytes))
	if err := decoder.Decode(&session); err != nil {
		return entity.Session{}, err
	}

	return session, nil
}

// Set a session by token
func (p *RedisProvider) Set(token string, session entity.Session) error {
	var sessionBytes bytes.Buffer

	encoder := gob.NewEncoder(&sessionBytes)
	if err := encoder.Encode(session); err != nil {
		return err
	}

	sessionString := hex.EncodeToString(sessionBytes.Bytes())

	if _, err := p.Connection.Do("SETEX", token, Duration.Seconds(), sessionString); err != nil {
		return err
	}

	return nil
}

// AutoMigrate the data connection
func (p *RedisProvider) AutoMigrate(clearExitstingData bool) error {
	if clearExitstingData {
		p.Connection.Do("FLUSHALL")
	}

	return nil
}
