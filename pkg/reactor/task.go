package reactor

import (
	"time"

	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/google/uuid"
)

// DatabaseTaskToEntity does what it says
func DatabaseTaskToEntity(dbTask entity.DatabaseTask) entity.Task {
	return entity.Task{
		CanBeCompletedEarly: dbTask.CanBeCompletedEarly,
		Completed:           dbTask.Completed,
		CreatedAt:           dbTask.CreatedAt,
		DeletedAt:           dbTask.DeletedAt,
		DueDate:             dbTask.DueDate,
		Name:                dbTask.Name,
		Notes:               dbTask.Notes,
		UserUUID:            dbTask.UserUUID,
		UUID:                dbTask.UUID,
	}
}

// DatabaseTasksToEntity does what it says
func DatabaseTasksToEntity(dbTasks []entity.DatabaseTask) []entity.Task {
	tasks := []entity.Task{}
	for _, dbTask := range dbTasks {
		tasks = append(tasks, DatabaseTaskToEntity(dbTask))
	}

	return tasks
}

// TaskCreateRequestToDatabaseTask does what it says
func TaskCreateRequestToDatabaseTask(request entity.TaskCreateRequest, loggedInUser entity.UserContext) entity.DatabaseTask {
	return entity.DatabaseTask{
		CanBeCompletedEarly: request.CanBeCompletedEarly,
		Completed:           false,
		CreatedAt:           time.Now(),
		DeletedAt:           nil,
		DueDate:             request.DueDate,
		Name:                request.Name,
		Notes:               request.Notes,
		UserUUID:            loggedInUser.User.UUID,
		UUID:                uuid.New().String(),
	}
}
