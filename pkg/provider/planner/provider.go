package plan

import "github.com/fuzzingbits/hub/pkg/entity"

// Provider is a Planner data provider
type Provider interface {
	// Get a planner by user and date
	Get(userUUID string, date string) (entity.Planner, error)
	// Save a planner by creating it or updating an existing one
	Save(plan entity.Planner) error
	// Delete an existing planner
	Delete(plan entity.Planner) error
	// GetAll planners for a given user
	GetAll(userUUID string) ([]entity.Planner, error)
}
