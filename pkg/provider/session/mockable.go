package session

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable is a Redis SessionProvider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// Get a session by token
func (p *Mockable) Get(token string) (string, error) {
	result, err := p.Provider.GetByID(token)
	if err != nil {
		if errors.Is(err, mockableprovider.ErrNotFound) {
			return "", ErrNotFound
		}

		return "", err
	}

	return result.(string), nil
}

// Set a session by token
func (p *Mockable) Set(token string, userUUID string) error {
	return p.Provider.Create(token, userUUID)
}
