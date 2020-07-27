package entity

// UserSession is all of the user data needed for the session
type UserSession struct {
	User     User         `json:"user"`
	Settings UserSettings `json:"userSettings"`
}

// User for Hub Users
type User struct {
	UUID      string `json:"uuid"`
	Username  string `json:"username"`
	Password  string `json:"-"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// DatabaseUser the the database structure for User objects
type DatabaseUser struct {
	ID        uint   `gorm:"primary_key"`
	UUID      string `gorm:"size:36;not null"`
	Username  string `gorm:"size:32;not null"`
	Password  string `gorm:"size:64;not null"`
	Email     string `gorm:"size:64;not null"`
	FirstName string `gorm:"size:64;not null"`
	LastName  string `gorm:"size:64;not null"`
}

// TableName for GORM
func (d DatabaseUser) TableName() string {
	return "user"
}

// UserSettings for a User
type UserSettings struct {
	ThemeColor string `json:"themeColor"`
}

// CreateUserRequest is the request for creating users
type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}