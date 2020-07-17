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

// UserSettings for a User
type UserSettings struct {
	ThemeColor string `json:"themeColor"`
}

// CreateUserRequest is the request for creating users
type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Username  string `json:"useranmae"`
	Password  string `json:"password"`
}
