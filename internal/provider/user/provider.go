package user

import "github.com/fuzzingbits/hub/internal/entity"

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
