package hub

import "github.com/fuzzingbits/hub/internal/entity"

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

	setupRequired := false
	if len(users) < 1 {
		setupRequired = true
	}

	return entity.ServerStatus{
		SetupRequired: setupRequired,
	}, nil
}
