package user

import (
	"database/sql"

	"github.com/fuzzingbits/hub/internal/entity"
)

// Provider is for working with User data
type Provider interface {
	GetUserByUUID(uuid string) (entity.User, error)
	GetAll() ([]entity.User, error)
	Save(user entity.User) (entity.User, error)
	Delete(user entity.User) error
	Create(user entity.User) (entity.User, error)
}

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	Database *sql.DB
}

// GetUserByUUID gets a User by UUID
func (d *DatabaseProvider) GetUserByUUID(uuid string) (entity.User, error) {
	return entity.User{}, nil
}

// GetAll Users
func (d *DatabaseProvider) GetAll() ([]entity.User, error) {
	return []entity.User{}, nil
}

// Save a User
func (d *DatabaseProvider) Save(user entity.User) (entity.User, error) {
	return user, nil
}

// Delete a User
func (d *DatabaseProvider) Delete(user entity.User) error {
	return nil
}

// Create a User
func (d *DatabaseProvider) Create(user entity.User) (entity.User, error) {
	return user, nil
}
