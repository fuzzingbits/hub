package routine

import (
	"github.com/fuzzingbits/hub/pkg/entity"
)

// Provider is for working with User data
type Provider interface {
	// GetByUUID gets a routine by UUID
	GetByUUID(uuid string) (entity.DatabaseRoutine, error)
	// GetAllByUser routines for a user
	GetAllByUser(userUUID string) ([]entity.DatabaseRoutine, error)
	// Update a routine
	Update(routine *entity.DatabaseRoutine) error
	// Delete a routine
	Delete(routine entity.DatabaseRoutine) error
	// Create a routine
	Create(routine *entity.DatabaseRoutine) error
}
