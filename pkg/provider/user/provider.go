package user

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// ErrNotFound is when a user can not be found by the provided UUID
var ErrNotFound = errors.New("User Not Found")

// Provider is for working with User data
type Provider interface {
	// GetByUUID gets a User by UUID
	GetByUUID(uuid string) (entity.DatabaseUser, error)
	// GetAll Users
	GetAll() ([]entity.DatabaseUser, error)
	// Update a User
	Update(user *entity.DatabaseUser) error
	// Delete a User
	Delete(user entity.DatabaseUser) error
	// Create a User
	Create(user *entity.DatabaseUser) error
}
