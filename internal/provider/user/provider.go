package user

import (
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
