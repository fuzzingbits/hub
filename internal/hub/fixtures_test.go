package hub

import (
	"errors"
	"testing"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestCreateFixturesSuccess(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)
	if err := s.CreateFixtures(); err != nil {
		t.Error(err)
	}

	c.UserProviderValue.Provider.CreateError = errors.New("foobar")
	if err := s.CreateFixtures(); err == nil {
		t.Error("there should have been an error")
	}
}
