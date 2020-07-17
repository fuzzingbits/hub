package user

// User the the database structure for User objects
type User struct {
	ID        uint   `gorm:"primary_key"`
	UUID      string `gorm:"size:36;not null"`
	Username  string `gorm:"size:32;not null"`
	Password  string `gorm:"size:64;not null"`
	Email     string `gorm:"size:64;not null"`
	FirstName string `gorm:"size:64;not null"`
	LastName  string `gorm:"size:64;not null"`
}

// Provider is for working with User data
type Provider interface {
	// GetByUUID gets a User by UUID
	GetByUUID(uuid string) (User, error)
	// GetAll Users
	GetAll() ([]User, error)
	// Update a User
	Update(user *User) error
	// Delete a User
	Delete(user User) error
	// Create a User
	Create(user *User) error
}
