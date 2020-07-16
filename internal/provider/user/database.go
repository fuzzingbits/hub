package user

import (
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/jinzhu/gorm"
)

type databaseUser struct {
	ID        uint   `gorm:"primary_key"`
	UUID      string `gorm:"size:64;not null"`
	FirstName string `gorm:"size:64;not null"`
	LastName  string `gorm:"size:64;not null"`
}

func (d databaseUser) TableName() string {
	return "user"
}

func databaseUserToEntity(dbUser databaseUser) entity.User {
	return entity.User{
		UUID:      dbUser.UUID,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
	}
}

func entityToDatabaseUser(entityUser entity.User, dbUser *databaseUser) {
	// Copy over values from the entityUser to an existing databaseUser
	dbUser.UUID = entityUser.UUID
	dbUser.FirstName = entityUser.FirstName
	dbUser.LastName = entityUser.LastName
}

// DatabaseProvider is a user.Provider the uses a database
type DatabaseProvider struct {
	Database *gorm.DB
}

// AutoMigrate the data connection
func (d *DatabaseProvider) AutoMigrate(clearExitstingData bool) error {
	// If devMode clear the table first
	if clearExitstingData {
		d.Database.DropTableIfExists(&databaseUser{})
	}

	// Always automigrate the table
	if err := d.Database.AutoMigrate(databaseUser{}).Error; err != nil {
		return err
	}

	return nil
}

// GetByUUID gets a User by UUID
func (d *DatabaseProvider) GetByUUID(uuid string) (entity.User, error) {
	dbUser, err := d.getByUUID(uuid)
	if err != nil {
		return entity.User{}, err
	}

	return databaseUserToEntity(dbUser), nil
}

// GetAll Users
func (d *DatabaseProvider) GetAll() ([]entity.User, error) {
	entityUsers := []entity.User{}
	dbUsers := []databaseUser{}
	if err := d.Database.Find(dbUsers).Error; err != nil {
		return nil, err
	}

	for _, dbUser := range dbUsers {
		entityUsers = append(entityUsers, databaseUserToEntity(dbUser))
	}

	return entityUsers, nil
}

// Save a User
func (d *DatabaseProvider) Save(user entity.User) (entity.User, error) {
	dbUser, err := d.getByUUID(user.UUID)
	if err != nil {
		return entity.User{}, err
	}

	entityToDatabaseUser(user, &dbUser)

	if err := d.Database.Save(&dbUser).Error; err != nil {
		return entity.User{}, err
	}

	return databaseUserToEntity(dbUser), nil
}

// Delete a User
func (d *DatabaseProvider) Delete(user entity.User) error {
	if err := d.Database.Where("`uuid` LIKE ?", user.UUID).Delete(databaseUser{}).Error; err != nil {
		return err
	}

	return nil
}

// Create a User
func (d *DatabaseProvider) Create(user entity.User) (entity.User, error) {
	dbUser := databaseUser{}
	entityToDatabaseUser(user, &dbUser)

	if err := d.Database.Create(&dbUser).Error; err != nil {
		return entity.User{}, err
	}

	user = databaseUserToEntity(dbUser)

	return user, nil
}

func (d *DatabaseProvider) getByUUID(uuid string) (databaseUser, error) {
	var dbUser databaseUser

	if err := d.Database.Where("`uuid` = ?", uuid).First(&dbUser).Error; err != nil {
		return databaseUser{}, err
	}

	return dbUser, nil
}
