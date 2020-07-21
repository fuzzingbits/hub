package usersettings

import "github.com/fuzzingbits/hub/pkg/entity"

// Provider is for working with UserSettings
type Provider interface {
	// GetByUUID gets a UserSettings by User.UUID
	GetByUUID(uuid string) (entity.UserSettings, error)
	// Save a UserSettings
	Save(uuid string, userSettings entity.UserSettings) error
	// Delete a UserSettings by UUID
	Delete(uuid string) error
}
