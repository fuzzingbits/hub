package usersettings

import (
	"sync"

	"github.com/fuzzingbits/hub/internal/entity"
)

type mockableEntry struct {
	UUID         string
	userSettings entity.UserSettings
}

// Mockable usersettings.Provider
type Mockable struct {
	Mutex          *sync.Mutex
	GetByUUIDError error
	SaveError      error
	DeleteError    error
	store          []mockableEntry
}

// GetByUUID a UserSettings
func (m *Mockable) GetByUUID(uuid string) (entity.UserSettings, error) {
	if m.GetByUUIDError != nil {
		return entity.UserSettings{}, m.GetByUUIDError
	}

	for _, entry := range m.store {
		if entry.UUID == uuid {
			return entry.userSettings, nil
		}
	}

	return entity.UserSettings{}, nil
}

// Save a UserSettings
func (m *Mockable) Save(uuid string, userSettings entity.UserSettings) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if m.SaveError != nil {
		return m.SaveError
	}

	for i, entry := range m.store {
		if entry.UUID == uuid {
			m.store[i].userSettings = userSettings
			return nil
		}
	}

	m.store = append(m.store, mockableEntry{UUID: uuid, userSettings: userSettings})

	return nil
}

// Delete a UserSettings
func (m *Mockable) Delete(uuid string) error {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	if m.DeleteError != nil {
		return m.DeleteError
	}

	for i, entry := range m.store {
		if entry.UUID == uuid {
			m.store = append(m.store[:i], m.store[i+1:]...)
			return nil
		}
	}

	return nil
}
