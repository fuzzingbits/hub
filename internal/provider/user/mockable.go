package user

import (
	"errors"
	"sync"

	"github.com/fuzzingbits/hub/internal/entity"
)

// ErrNotFound is when a user can not be found by the provided UUID
var ErrNotFound = errors.New("User Not Found")

// Mockable user.Provider
type Mockable struct {
	users       []entity.User
	Mutex       *sync.Mutex
	GetAllError error
	CreateError error
	SaveError   error
	DeleteError error
}

// GetUserByUUID gets a user by UUID
func (m *Mockable) GetUserByUUID(uuid string) (entity.User, error) {
	for _, user := range m.users {
		if user.UUID == uuid {
			return user, nil
		}
	}

	return entity.User{}, ErrNotFound
}

// GetAll users
func (m *Mockable) GetAll() ([]entity.User, error) {
	if m.GetAllError != nil {
		return nil, m.GetAllError
	}

	return m.users, nil
}

// Save a user
func (m *Mockable) Save(user entity.User) (entity.User, error) {
	if m.SaveError != nil {
		return entity.User{}, m.SaveError
	}

	m.Delete(user)
	m.Create(user)

	return user, nil
}

// Delete a user
func (m *Mockable) Delete(user entity.User) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if m.DeleteError != nil {
		return m.DeleteError
	}

	for i, storedUser := range m.users {
		if user.UUID == storedUser.UUID {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

// Create a user
func (m *Mockable) Create(user entity.User) (entity.User, error) {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if m.CreateError != nil {
		return user, m.CreateError
	}

	m.users = append(m.users, user)

	return user, nil
}
