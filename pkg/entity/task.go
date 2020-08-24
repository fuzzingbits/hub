package entity

import "time"

// Task is a todo item
type Task struct {
	UUID                string     `json:"uuid"`
	UserUUID            string     `json:"userUUID"`
	Name                string     `json:"name"`
	Note                string     `json:"note"`
	DueDate             string     `json:"dueDate"`
	Completed           bool       `json:"completed"`
	CreatedAt           time.Time  `json:"createdAt"`
	DeletedAt           *time.Time `json:"deletedAt"`
	CanBeCompletedEarly bool       `json:"canBeCompletedEarly"`
}

// DatabaseTask is a database task
type DatabaseTask struct {
	ID                  uint       `gorm:"primary_key"`
	UUID                string     `gorm:"size:36;not null"`
	UserUUID            string     `gorm:"size:36;not null"`
	Name                string     `gorm:"size:128;not null"`
	Note                string     `gorm:"type:longtext;not null"`
	DueDate             string     `gorm:"size:10;not null"`
	Completed           bool       `gorm:"not null"`
	CreatedAt           time.Time  `gorm:"not null"`
	DeletedAt           *time.Time `gorm:"null"`
	CanBeCompletedEarly bool       `gorm:"not null"`
}

// TableName for GORM
func (d DatabaseTask) TableName() string {
	return "task"
}

// TaskCreateRequest is a create task request
type TaskCreateRequest struct {
	Name                string `json:"name"`
	Note                string `json:"note"`
	DueDate             string `json:"dueDate"`
	Completed           bool   `json:"completed"`
	CanBeCompletedEarly bool   `json:"canBeCompletedEarly"`
}
