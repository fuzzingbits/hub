package hub

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// ErrServerAlreadySetup is when the server is already setup
var ErrServerAlreadySetup = errors.New("Server Already Setup")

// GetServerStatus gets the status for the server
func (s *Service) GetServerStatus() (entity.ServerStatus, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.ServerStatus{}, err
	}

	users, err := userProvider.GetAll()
	if err != nil {
		return entity.ServerStatus{}, err
	}

	setupRequired := len(users) < 1

	return entity.ServerStatus{
		SetupRequired: setupRequired,
	}, nil
}

// SetupServer sets up the server
func (s *Service) SetupServer(createUserRequest entity.CreateUserRequest) (entity.Session, error) {
	serverStatus, err := s.GetServerStatus()
	if err != nil {
		return entity.Session{}, err
	}

	if !serverStatus.SetupRequired {
		return entity.Session{}, ErrServerAlreadySetup
	}

	_, createUserErr := s.CreateUser(createUserRequest)
	if createUserErr != nil {
		return entity.Session{}, createUserErr
	}

	return s.Login(entity.UserLoginRequest{
		Email:    createUserRequest.Email,
		Password: createUserRequest.Password,
	})
}
