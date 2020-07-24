package user

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
		d.Database.DropTableIfExists(&entity.DatabaseUser{})
	}

	// Always automigrate the table
	if err := d.Database.AutoMigrate(entity.DatabaseUser{}).Error; err != nil {
		return err
	}

	return nil
}

// GetByUUID gets a User by UUID
func (d *DatabaseProvider) GetByUUID(uuid string) (entity.DatabaseUser, error) {
	var dbUser entity.DatabaseUser

	if err := d.Database.Where("`uuid` = ?", uuid).First(&dbUser).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return entity.DatabaseUser{}, ErrNotFound
		}

		return entity.DatabaseUser{}, err
	}

	return dbUser, nil
}

// GetAll Users
func (d *DatabaseProvider) GetAll() ([]entity.DatabaseUser, error) {
	dbUsers := []entity.DatabaseUser{}
	if err := d.Database.Find(dbUsers).Error; err != nil {
		return nil, err
	}

	return dbUsers, nil
}

// Update a User
func (d *DatabaseProvider) Update(dbUser *entity.DatabaseUser) error {
	if err := d.Database.Save(dbUser).Error; err != nil {
		return err
	}

	return nil
}

// Delete a User
func (d *DatabaseProvider) Delete(user entity.DatabaseUser) error {
	if err := d.Database.Where("`uuid` LIKE ?", user.UUID).Delete(entity.DatabaseUser{}).Error; err != nil {
		return err
	}

	return nil
}

// Create a User
func (d *DatabaseProvider) Create(dbUser *entity.DatabaseUser) error {
	if err := d.Database.Create(&dbUser).Error; err != nil {
		return err
	}

	return nil
}
