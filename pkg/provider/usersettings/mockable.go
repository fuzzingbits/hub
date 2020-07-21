package usersettings

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable usersettings.Provider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// GetByUUID gets a UserSettings by User.UUID
func (m *Mockable) GetByUUID(uuid string) (entity.UserSettings, error) {
	item, err := m.Provider.GetByID(uuid)
	if err != nil {
		return entity.UserSettings{}, err
	}

	userSettings, _ := item.(entity.UserSettings)

	return userSettings, nil
}

// Save a UserSettings
func (m *Mockable) Save(uuid string, userSettings entity.UserSettings) error {
	return m.Provider.Update(uuid, userSettings)
}

// Delete a UserSettings by UUID
func (m *Mockable) Delete(uuid string) error {
	return m.Provider.Delete(uuid)
}
