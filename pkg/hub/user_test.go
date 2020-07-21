package hub

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
)

func TestGetCurrentSession(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	// Create Fixture User
	targetUserSession, _ := s.CreateUser(
		entity.CreateUserRequest{
			FirstName: "Testy",
			LastName:  "McTestPants",
			Username:  "testy",
			Email:     "testy@example.com",
			Password:  "Password123",
		},
	)

	// Create Fake Request
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("UUID", targetUserSession.User.UUID)

	{ // Test Success!
		session, err := s.GetCurrentSession(req)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(session, targetUserSession) {
			t.Errorf(
				"[Session Did Not Match] returned: %s expected: %s",
				session,
				targetUserSession,
			)
		}
	}

	{ // Test Get By UUID Failure
		c.UserSettingsProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}

	{ // Test Get By UUID Failure
		c.UserProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}

	{ // Test Get By UUID Failure
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}

	{ // Test Get By UUID Failure
		c.UserProviderError = errors.New("foobar")
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}
}

func TestCreateUser(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	createUserRequest := entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	}

	if _, err := s.CreateUser(createUserRequest); err != nil {
		t.Error(err)
	}

	c.UserSettingsProviderValue.Provider.UpdateError = errors.New("foobar")
	if _, err := s.CreateUser(createUserRequest); err == nil {
		t.Error("there should have been an error")
	}

	c.UserSettingsProviderError = errors.New("foobar")
	if _, err := s.CreateUser(createUserRequest); err == nil {
		t.Error("there should have been an error")
	}

	c.UserProviderError = errors.New("foobar")
	if _, err := s.CreateUser(createUserRequest); err == nil {
		t.Error("there should have been an error")
	}
}
