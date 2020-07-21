package hub

import "github.com/fuzzingbits/hub/pkg/entity"

// CreateFixtures creates fixtures
func (s *Service) CreateFixtures() error {
	if _, err := s.CreateUser(entity.CreateUserRequest{
		FirstName: "Aaron",
		LastName:  "Ellington",
		Username:  "aaron",
		Email:     "aaron@example.com",
		Password:  "hub100",
	}); err != nil {
		return err
	}

	return nil
}
