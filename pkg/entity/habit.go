package entity

// HabitStore is a storage document for habits
type HabitStore struct {
	UserUUID string  `json:"userUUID"`
	Habits   []Habit `json:"habits"`
}

// Habit is an action that should be preformed on a schedule
type Habit struct {
	Name      string `json:"name"`
	Sunday    bool   `json:"sunday"`
	Monday    bool   `json:"monday"`
	Tuesday   bool   `json:"tuesday"`
	Wednesday bool   `json:"wednesday"`
	Thursday  bool   `json:"thursday"`
	Friday    bool   `json:"friday"`
	Saturday  bool   `json:"saturday"`
}
