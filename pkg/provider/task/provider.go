package task

import "github.com/fuzzingbits/hub/pkg/entity"

// Provider is for working with Task data
type Provider interface {
	// Create a task
	Create(task *entity.DatabaseTask) error
	// Update a task
	Update(task *entity.DatabaseTask) error
	// Delete a task
	Delete(task entity.DatabaseTask) error
	// GetByUUID gets a task by UUID
	GetByUUID(uuid string) (entity.DatabaseTask, error)
	// GetAllActive tasks
	GetAllActive(userUUID string) ([]entity.DatabaseTask, error)
	// GetAllDeleted tasks
	GetAllDeleted(userUUID string) ([]entity.DatabaseTask, error)
	// GetAllByDueDate tasks
	GetAllByDueDate(userUUID string, dueDate string) ([]entity.DatabaseTask, error)
	// GetAllDueNow tasks
	GetAllDueNow(userUUID string) ([]entity.DatabaseTask, error)
	// GetAllDueSoon tasks
	GetAllDueSoon(userUUID string) ([]entity.DatabaseTask, error)
	// GetAllStale tasks
	GetAllStale(userUUID string) ([]entity.DatabaseTask, error)
}
