package entity

// Routine is an action that you would like to regularly complete
type Routine struct {
	Name     string          `json:"name"`
	Enabled  bool            `json:"enabled"`
	Schedule RoutineSchedule `json:"routineSchedule"`
}

// RoutineSchedule is the schedule for a given routine
type RoutineSchedule struct {
	Sunday    bool `json:"sunday"`
	Monday    bool `json:"monday"`
	Tuesday   bool `json:"tuesday"`
	Wednesday bool `json:"wednesday"`
	Thursday  bool `json:"thursday"`
	Friday    bool `json:"friday"`
	Saturday  bool `json:"saturday"`
}

// DatabaseRoutine is the databae representation of a routine
type DatabaseRoutine struct {
	ID        uint   `gorm:"primary_key"`
	UUID      string `gorm:"size:36;not null"`
	UserUUID  string `gorm:"size:36;not null"`
	Name      string `gorm:"size:128;not null"`
	Enabled   bool   `gorm:"not null"`
	Sunday    bool   `gorm:"not null"`
	Monday    bool   `gorm:"not null"`
	Tuesday   bool   `gorm:"not null"`
	Wednesday bool   `gorm:"not null"`
	Thursday  bool   `gorm:"not null"`
	Friday    bool   `gorm:"not null"`
	Saturday  bool   `gorm:"not null"`
}
