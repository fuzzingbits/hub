package user

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable user.Provider
type Mockable struct {
	Provider        *mockableprovider.Provider
	GetByEmailError error
}

// GetByUUID gets a user by UUID
func (m *Mockable) GetByUUID(uuid string) (entity.DatabaseUser, error) {
	item, err := m.Provider.GetByID(uuid)
	if err != nil {
		return entity.DatabaseUser{}, err
	}

	user, _ := item.(entity.DatabaseUser)

	return user, nil
}

// GetByEmail gets a user by email
func (m *Mockable) GetByEmail(email string) (entity.DatabaseUser, error) {
	if m.GetByEmailError != nil {
		return entity.DatabaseUser{}, m.GetByEmailError
	}

	items, err := m.Provider.GetAll()
	if err != nil {
		return entity.DatabaseUser{}, err
	}

	for _, item := range items {
		user, _ := item.(entity.DatabaseUser)
		if user.Email == email {
			return user, nil
		}
	}

	return entity.DatabaseUser{}, entity.ErrRecordNotFound
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
