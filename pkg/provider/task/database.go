package task

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/jinzhu/gorm"
)

// DatabaseProvider is a task.Provider the uses a database
type DatabaseProvider struct {
	Database *gorm.DB
}

// AutoMigrate the data connection
func (d *DatabaseProvider) AutoMigrate(clearExitstingData bool) error {
	// If devMode clear the table first
	if clearExitstingData {
		d.Database.DropTableIfExists(&entity.DatabaseTask{})
	}

	// Always automigrate the table
	if err := d.Database.AutoMigrate(entity.DatabaseTask{}).Error; err != nil {
		return err
	}

	return nil
}

// Create a task
func (d *DatabaseProvider) Create(task *entity.DatabaseTask) error {
	if err := d.Database.Create(&task).Error; err != nil {
		return err
	}

	return nil
}

// Update a task
func (d *DatabaseProvider) Update(task *entity.DatabaseTask) error {
	if err := d.Database.Save(task).Error; err != nil {
		return err
	}

	return nil
}

// Delete a task
func (d *DatabaseProvider) Delete(task entity.DatabaseTask) error {
	if err := d.Database.Where("`uuid` = ?", task.UUID).Delete(entity.DatabaseTask{}).Error; err != nil {
		return err
	}

	return nil
}

// GetByUUID gets a task by UUID
func (d *DatabaseProvider) GetByUUID(uuid string) (entity.DatabaseTask, error) {
	return entity.DatabaseTask{}, nil
}

// GetAllActive tasks
func (d *DatabaseProvider) GetAllActive(userUUID string) ([]entity.DatabaseTask, error) {
	return d.getAllHelper(
		`
			SELECT *
			FROM task
			WHERE
				user_uuid = ?
				AND deleted_at IS NULL
		`,
		userUUID,
	)
}

// GetAllDeleted tasks
func (d *DatabaseProvider) GetAllDeleted(userUUID string) ([]entity.DatabaseTask, error) {
	return d.getAllHelper(
		`
			SELECT *
			FROM task
			WHERE
				user_uuid = ?
				AND deleted_at IS NOT NULL
		`,
		userUUID,
	)
}

// GetAllByDueDate tasks
func (d *DatabaseProvider) GetAllByDueDate(userUUID string, dueDate string) ([]entity.DatabaseTask, error) {
	return d.getAllHelper(
		`
			SELECT *
			FROM task
			WHERE
				user_uuid = ?
				AND deleted_at IS NULL
				AND due_date = ?
		`,
		userUUID,
		dueDate,
	)
}

// GetAllDueNow tasks
func (d *DatabaseProvider) GetAllDueNow(userUUID string) ([]entity.DatabaseTask, error) {
	return d.getAllHelper(
		`
			SELECT *
			FROM task
			WHERE
				user_uuid = ?
				AND deleted_at IS NULL
				AND DATE(due_date) <= DATE(CURRENT_TIMESTAMP)
		`,
		userUUID,
	)
}

// GetAllDueSoon tasks
func (d *DatabaseProvider) GetAllDueSoon(userUUID string) ([]entity.DatabaseTask, error) {
	// TODO: Complete this function
	return []entity.DatabaseTask{}, nil
}

// GetAllStale tasks
func (d *DatabaseProvider) GetAllStale(userUUID string) ([]entity.DatabaseTask, error) {
	// TODO: Complete this function
	return []entity.DatabaseTask{}, nil
}

func (d *DatabaseProvider) getAllHelper(sql string, values ...interface{}) ([]entity.DatabaseTask, error) {
	tasks := []entity.DatabaseTask{}
	if err := d.Database.Raw(sql, values...).Scan(tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}
