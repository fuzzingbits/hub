package session

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable is a Redis SessionProvider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// Get a session by token
func (p *Mockable) Get(token string) (entity.UserSession, error) {
	result, err := p.Provider.GetByID(token)
	if err != nil {
		return entity.UserSession{}, err
	}

	return result.(entity.UserSession), nil
}

// Set a session by token
func (p *Mockable) Set(token string, session entity.UserSession) error {
	return p.Provider.Create(token, session)
}
