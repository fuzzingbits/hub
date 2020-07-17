package hub

// CreateFixtures creates fixtures
func (s *Service) CreateFixtures() error {
	if _, err := s.CreateUser("313efbe9-173b-4a1b-9a5b-7b69d95a66b9", "Aaron", "Ellington"); err != nil {
		return err
	}

	return nil
}
