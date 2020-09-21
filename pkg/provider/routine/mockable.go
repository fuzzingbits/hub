package routine

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable user.Provider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// GetByUUID gets a routine by UUID
func (m *Mockable) GetByUUID(uuid string) (entity.DatabaseRoutine, error) {
	item, err := m.Provider.GetByID(uuid)
	if err != nil {
		return entity.DatabaseRoutine{}, err
	}

	user, _ := item.(entity.DatabaseRoutine)

	return user, nil
}

// GetAllByUser routines for a user
func (m *Mockable) GetAllByUser(userUUID string) ([]entity.DatabaseRoutine, error) {
	items, err := m.Provider.GetAll()
	if err != nil {
		return nil, err
	}

	users := []entity.DatabaseRoutine{}
	for _, item := range items {
		user, _ := item.(entity.DatabaseRoutine)
		users = append(users, user)
	}

	return users, nil
}

// Update a routine
func (m *Mockable) Update(user *entity.DatabaseRoutine) error {
	return m.Provider.Update(user.UUID, *user)
}

// Delete a routine
func (m *Mockable) Delete(user entity.DatabaseRoutine) error {
	return m.Provider.Delete(user.UUID)
}

// Create a routine
func (m *Mockable) Create(user *entity.DatabaseRoutine) error {
	return m.Provider.Create(user.UUID, *user)
}
