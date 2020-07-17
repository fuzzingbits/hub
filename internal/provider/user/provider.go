package user

import (
	"github.com/fuzzingbits/hub/internal/entity"
)

// Provider is for working with User data
type Provider interface {
	// GetByUUID gets a User by UUID
	GetByUUID(uuid string) (entity.User, error)
	// GetAll Users
	GetAll() ([]entity.User, error)
	// Update a User
	Update(user entity.User) error
	// Delete a User
	Delete(user entity.User) error
	// Create a User
	Create(user entity.User) error
}
