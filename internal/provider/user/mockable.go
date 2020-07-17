package user

import (
	"errors"

	"github.com/fuzzingbits/hub/internal/forge/mockableprovider"
)

// ErrNotFound is when a user can not be found by the provided UUID
var ErrNotFound = errors.New("User Not Found")

// Mockable user.Provider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// GetByUUID gets a user by UUID
func (m *Mockable) GetByUUID(uuid string) (User, error) {
	item, err := m.Provider.GetByID(uuid)
	if err != nil {
		return User{}, err
	}

	user, _ := item.(User)

	return user, nil
}

// GetAll Users
func (m *Mockable) GetAll() ([]User, error) {
	items, err := m.Provider.GetAll()
	if err != nil {
		return nil, err
	}

	users := []User{}
	for _, item := range items {
		user, _ := item.(User)
		users = append(users, user)
	}

	return users, nil
}

// Update a User
func (m *Mockable) Update(user *User) error {
	return m.Provider.Update(user.UUID, *user)
}

// Delete a User
func (m *Mockable) Delete(user User) error {
	return m.Provider.Delete(user.UUID)
}

// Create a User
func (m *Mockable) Create(user *User) error {
	return m.Provider.Create(user.UUID, *user)
}
