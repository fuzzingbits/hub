package entity

import "time"

// Planner for a specified user and day
type Planner struct {
	UserUUID        string         `json:"userUUID"`
	Date            string         `json:"date"`
	Priorities      []string       `json:"priorities"`
	Accomplishments []string       `json:"accomplishments"`
	TasksToday      []PlannerTask  `json:"tasksToday"`
	TasksTomorrow   []PlannerTask  `json:"tasksTomorrow"`
	Schedule        []PlannerEvent `json:"schedule"`
}

// PlannerEvent that exists in a schedule
type PlannerEvent struct {
	Name  string    `json:"name"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Color string    `json:"color"`
}

// PlannerTask as part of a plan
type PlannerTask struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
