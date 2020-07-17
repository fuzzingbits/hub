package hub

import (
	"errors"
	"reflect"
	"testing"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestGetServerStatus(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	{ // Success! on new server
		status, err := s.GetServerStatus()
		if err != nil {
			t.Error(err)
		}

		targetStatus := entity.ServerStatus{
			SetupRequired: true,
		}

		if !reflect.DeepEqual(status, targetStatus) {
			t.Errorf(
				"[ServerStatus Did Not Match] returned: %+v expected: %+v",
				status,
				targetStatus,
			)
		}
	}

	{ // Success! on server with a user already
		s.CreateUser(entity.CreateUserRequest{
			FirstName: "Testy",
			LastName:  "McTestPants",
			Username:  "testy",
			Email:     "testy@example.com",
			Password:  "Password123",
		})

		status, err := s.GetServerStatus()
		if err != nil {
			t.Error(err)
		}

		targetStatus := entity.ServerStatus{
			SetupRequired: false,
		}

		if !reflect.DeepEqual(status, targetStatus) {
			t.Errorf(
				"[ServerStatus Did Not Match] returned: %+v expected: %+v",
				status,
				targetStatus,
			)
		}
	}

	{ // Failed, could not get users
		c.UserProviderValue.Provider.GetAllError = errors.New("foobar")
		if _, err := s.GetServerStatus(); err == nil {
			t.Error("there should have been an error")
		}
	}

	{ // Failed, could not get userProvider
		c.UserProviderError = errors.New("foobar")
		if _, err := s.GetServerStatus(); err == nil {
			t.Error("there should have been an error")
		}
	}
}
