package routine

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/jinzhu/gorm"
)

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	Database *gorm.DB
}

// AutoMigrate the data connection
func (d *DatabaseProvider) AutoMigrate(clearExitstingData bool) error {
	// If devMode clear the table first
	if clearExitstingData {
		d.Database.DropTableIfExists(&entity.DatabaseRoutine{})
	}

	// Always automigrate the table
	if err := d.Database.AutoMigrate(entity.DatabaseRoutine{}).Error; err != nil {
		return err
	}

	return nil
}

// GetByUUID gets a routine by UUID
func (d *DatabaseProvider) GetByUUID(uuid string) (entity.DatabaseRoutine, error) {
	var dbRoutine entity.DatabaseRoutine

	if err := d.Database.Where("`uuid` = ?", uuid).First(&dbRoutine).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity.DatabaseRoutine{}, entity.ErrRecordNotFound
		}

		return entity.DatabaseRoutine{}, err
	}

	return dbRoutine, nil
}

// GetAllByUser routines for a user
func (d *DatabaseProvider) GetAllByUser(userUUID string) ([]entity.DatabaseRoutine, error) {
	dbRoutines := []entity.DatabaseRoutine{}
	if err := d.Database.Find(&dbRoutines).Error; err != nil {
		return nil, err
	}

	return dbRoutines, nil
}

// Update a routine
func (d *DatabaseProvider) Update(dbRoutine *entity.DatabaseRoutine) error {
	if err := d.Database.Save(dbRoutine).Error; err != nil {
		return err
	}

	return nil
}

// Delete a routine
func (d *DatabaseProvider) Delete(user entity.DatabaseRoutine) error {
	if err := d.Database.Where("`uuid` LIKE ?", user.UUID).Delete(entity.DatabaseRoutine{}).Error; err != nil {
		return err
	}

	return nil
}

// Create a routine
func (d *DatabaseProvider) Create(dbRoutine *entity.DatabaseRoutine) error {
	if err := d.Database.Create(&dbRoutine).Error; err != nil {
		return err
	}

	return nil
}
