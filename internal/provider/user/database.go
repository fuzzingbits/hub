package user

import (
	"github.com/jinzhu/gorm"
)

// TableName for GORM
func (d User) TableName() string {
	return "user"
}

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	Database *gorm.DB
}

// AutoMigrate the data connection
func (d *DatabaseProvider) AutoMigrate(clearExitstingData bool) error {
	// If devMode clear the table first
	if clearExitstingData {
		d.Database.DropTableIfExists(&User{})
	}

	// Always automigrate the table
	if err := d.Database.AutoMigrate(User{}).Error; err != nil {
		return err
	}

	return nil
}

// GetByUUID gets a User by UUID
func (d *DatabaseProvider) GetByUUID(uuid string) (User, error) {
	return d.getByUUID(uuid)
}

// GetAll Users
func (d *DatabaseProvider) GetAll() ([]User, error) {
	dbUsers := []User{}
	if err := d.Database.Find(dbUsers).Error; err != nil {
		return nil, err
	}

	return dbUsers, nil
}

// Update a User
func (d *DatabaseProvider) Update(dbUser *User) error {
	if err := d.Database.Save(dbUser).Error; err != nil {
		return err
	}

	return nil
}

// Delete a User
func (d *DatabaseProvider) Delete(user User) error {
	if err := d.Database.Where("`uuid` LIKE ?", user.UUID).Delete(User{}).Error; err != nil {
		return err
	}

	return nil
}

// Create a User
func (d *DatabaseProvider) Create(dbUser *User) error {

	if err := d.Database.Create(&dbUser).Error; err != nil {
		return err
	}

	return nil
}

func (d *DatabaseProvider) getByUUID(uuid string) (User, error) {
	var dbUser User

	if err := d.Database.Where("`uuid` = ?", uuid).First(&dbUser).Error; err != nil {
		return User{}, err
	}

	return dbUser, nil
}
