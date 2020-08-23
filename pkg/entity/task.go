package entity

// Task is a todo item
type Task struct {
	UUID      string `json:"uuid"`
	UserUUID  string `json:"userUUID"`
	Category  string `json:"category"`
	Name      string `json:"name"`
	Note      string `json:"note"`
	DueDate   string `json:"dueDate"`
	Completed bool   `json:"completed"`
	// TODO: Come up with a better name for this
	AllowedToBeCompletedEarly bool `json:"allowedToBeCompletedEarly"`
}
