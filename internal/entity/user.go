package entity

// UserSession is all of the user data needed for the session
type UserSession struct {
	User     User         `json:"user"`
	Settings UserSettings `json:"userSettings"`
}

// User for Hub Users
type User struct {
	UUID      string `json:"uuid"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// UserSettings for a User
type UserSettings struct {
	ThemeColor string `json:"themeColor"`
}
