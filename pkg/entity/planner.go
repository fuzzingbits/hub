package entity

import "time"

// Planner for a specified user and day
type Planner struct {
	UserUUID        string         `json:"userUUID"`
	Date            string         `json:"date"`
	Updated         time.Time      `json:"updated"`
	Created         time.Time      `json:"created"`
	Priorities      []string       `json:"priorities"`
	Accomplishments []string       `json:"accomplishments"`
	TasksToday      []PlannerTask  `json:"tasksToday"`
	TasksTomorrow   []PlannerTask  `json:"tasksTomorrow"`
	Schedule        []PlannerEvent `json:"schedule"`
}

// PlannerEvent that exists in a schedule
type PlannerEvent struct {
	Value string    `json:"value"`
	End   time.Time `json:"end"`
	Start time.Time `json:"start"`
	Color string    `json:"color"`
}

// PlannerTask as part of a plan
type PlannerTask struct {
	Value     string `json:"value"`
	Completed bool   `json:"completed"`
}
