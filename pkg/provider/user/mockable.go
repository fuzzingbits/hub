package user

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable user.Provider
type Mockable struct {
	Provider           *mockableprovider.Provider
	GetByUsernameError error
}

// GetByUUID gets a user by UUID
func (m *Mockable) GetByUUID(uuid string) (entity.DatabaseUser, error) {
	item, err := m.Provider.GetByID(uuid)
	if err != nil {
		if errors.Is(err, mockableprovider.ErrNotFound) {
			return entity.DatabaseUser{}, ErrNotFound
		}

		return entity.DatabaseUser{}, err
	}

	user, _ := item.(entity.DatabaseUser)

	return user, nil
}

// GetByUsername gets a user by username
func (m *Mockable) GetByUsername(username string) (entity.DatabaseUser, error) {
	if m.GetByUsernameError != nil {
		return entity.DatabaseUser{}, m.GetByUsernameError
	}

	items, err := m.Provider.GetAll()
	if err != nil {
		return entity.DatabaseUser{}, err
	}

	for _, item := range items {
		user, _ := item.(entity.DatabaseUser)
		if user.Username == username {
			return user, nil
		}
	}

	return entity.DatabaseUser{}, ErrNotFound
}

// GetAll Users
func (m *Mockable) GetAll() ([]entity.DatabaseUser, error) {
	items, err := m.Provider.GetAll()
	if err != nil {
		return nil, err
	}

	users := []entity.DatabaseUser{}
	for _, item := range items {
		user, _ := item.(entity.DatabaseUser)
		users = append(users, user)
	}

	return users, nil
}

// Update a User
func (m *Mockable) Update(user *entity.DatabaseUser) error {
	return m.Provider.Update(user.UUID, *user)
}

// Delete a User
func (m *Mockable) Delete(user entity.DatabaseUser) error {
	return m.Provider.Delete(user.UUID)
}

// Create a User
func (m *Mockable) Create(user *entity.DatabaseUser) error {
	return m.Provider.Create(user.UUID, *user)
}
