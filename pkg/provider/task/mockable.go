package task

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/mockableprovider"
)

// Mockable user.Provider
type Mockable struct {
	Provider *mockableprovider.Provider
}

// Create a task
func (m *Mockable) Create(task *entity.DatabaseTask) error {
	return m.Provider.Create(task.UUID, *task)
}

// Update a task
func (m *Mockable) Update(task *entity.DatabaseTask) error {
	return m.Provider.Update(task.UUID, *task)
}

// Delete a task
func (m *Mockable) Delete(task entity.DatabaseTask) error {
	return m.Provider.Delete(task.UUID)
}

// GetByUUID gets a task by UUID
func (m *Mockable) GetByUUID(uuid string) (entity.DatabaseTask, error) {
	item, err := m.Provider.GetByID(uuid)
	if err != nil {
		return entity.DatabaseTask{}, err
	}

	task, _ := item.(entity.DatabaseTask)

	return task, nil
}

// GetAllActive tasks
func (m *Mockable) GetAllActive(userUUID string) ([]entity.DatabaseTask, error) {
	return []entity.DatabaseTask{}, nil
}

// GetAllDeleted tasks
func (m *Mockable) GetAllDeleted(userUUID string) ([]entity.DatabaseTask, error) {
	return []entity.DatabaseTask{}, nil
}

// GetAllByDueDate tasks
func (m *Mockable) GetAllByDueDate(userUUID string, dueDate string) ([]entity.DatabaseTask, error) {
	return []entity.DatabaseTask{}, nil
}

// GetAllDueNow tasks
func (m *Mockable) GetAllDueNow(userUUID string) ([]entity.DatabaseTask, error) {
	return []entity.DatabaseTask{}, nil
}

// GetAllDueSoon tasks
func (m *Mockable) GetAllDueSoon(userUUID string) ([]entity.DatabaseTask, error) {
	return []entity.DatabaseTask{}, nil
}

// GetAllStale tasks
func (m *Mockable) GetAllStale(userUUID string) ([]entity.DatabaseTask, error) {
	return []entity.DatabaseTask{}, nil
}
