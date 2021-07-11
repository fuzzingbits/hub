package hub

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/reactor"
)

// CreateTask creates a task
func (s *Service) CreateTask(loggedInUser entity.UserContext, request entity.TaskCreateRequest) (entity.Task, error) {
	taskProvider, err := s.container.TaskProvider()
	if err != nil {
		return entity.Task{}, err
	}

	dbTask := reactor.TaskCreateRequestToDatabaseTask(request, loggedInUser)

	if err := taskProvider.Create(&dbTask); err != nil {
		return entity.Task{}, err
	}

	return entity.Task{}, nil
}

// GetAllActiveTasks for a given user
func (s *Service) GetAllActiveTasks(loggedInUser entity.UserContext) ([]entity.Task, error) {
	taskProvider, err := s.container.TaskProvider()
	if err != nil {
		return []entity.Task{}, err
	}

	dbTasks, err := taskProvider.GetAllActive(loggedInUser.User.UUID)
	if err != nil {
		return []entity.Task{}, err
	}

	tasks := reactor.DatabaseTasksToEntity(dbTasks)

	return tasks, nil
}
