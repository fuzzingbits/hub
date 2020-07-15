package usersettings

import "github.com/fuzzingbits/hub/internal/entity"

// Provider is for working with UserSettings data
type Provider interface {
	GetByUUID(uuid string) (entity.UserSettings, error)
	Save(uuid string, userSettings entity.UserSettings) error
	Delete(uuid string) error
}
