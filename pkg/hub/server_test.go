package hub

import (
	"errors"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
)

func TestGetServerStatus(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	{ // Success
		if _, err := s.GetServerStatus(); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		c.UserProviderValue.Provider.GetAllError = errors.New("foobar")
		if _, err := s.GetServerStatus(); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.GetServerStatus(); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}

func TestSetupServer(t *testing.T) {
	{ // Success and already setup
		c := container.NewMockable()
		s := NewService(&hubconfig.Config{}, c)

		if _, err := s.SetupServer(standardTestUserCreateRequest); err != nil {
			t.Error(err)
		}

		if _, err := s.SetupServer(standardTestUserCreateRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	{ // Error
		c.UserProviderValue.Provider.CreateError = errors.New("foobar")
		if _, err := s.SetupServer(standardTestUserCreateRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.SetupServer(standardTestUserCreateRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}
