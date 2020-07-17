package hub

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestGetCurrentSession(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	targetUserSession := entity.UserSession{
		User: entity.User{
			UUID:      "313efbe9-173b-4a1b-9b5b-7b69d95a66b9",
			FirstName: "Testy",
			LastName:  "McTestPants",
		},
		Settings: entity.UserSettings{
			ThemeColor: "tomato",
		},
	}

	// Create Fixture User
	s.CreateUser(
		targetUserSession.User.UUID,
		targetUserSession.User.FirstName,
		targetUserSession.User.LastName,
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
		c.UserSettingsProviderValue.GetByUUIDError = errors.New("foobar")
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}

	{ // Test Get By UUID Failure
		c.UserProviderValue.GetByUUIDError = errors.New("foobar")
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

	if _, err := s.CreateUser(
		"fake-uuid",
		"fake-first-name",
		"fake-last-name",
	); err != nil {
		t.Error(err)
	}

	c.UserSettingsProviderValue.SaveError = errors.New("foobar")
	if _, err := s.CreateUser(
		"fake",
		"fake-first-name",
		"fake-last-name",
	); err == nil {
		t.Error("there should have been an error")
	}

	c.UserSettingsProviderError = errors.New("foobar")
	if _, err := s.CreateUser(
		"fake",
		"fake-first-name",
		"fake-last-name",
	); err == nil {
		t.Error("there should have been an error")
	}

	c.UserProviderError = errors.New("foobar")
	if _, err := s.CreateUser(
		"fake",
		"fake-first-name",
		"fake-last-name",
	); err == nil {
		t.Error("there should have been an error")
	}
}
