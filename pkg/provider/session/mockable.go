package session

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable is a Redis SessionProvider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// Get a session by token
func (p *Mockable) Get(token string) (entity.Session, error) {
	result, err := p.Provider.GetByID(token)
	if err != nil {
		if errors.Is(err, mockableprovider.ErrNotFound) {
			return entity.Session{}, ErrNotFound
		}

		return entity.Session{}, err
	}

	return result.(entity.Session), nil
}

// Set a session by token
func (p *Mockable) Set(token string, session entity.Session) error {
	return p.Provider.Create(token, session)
}
