package hub

import (
	"errors"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
)

var standardTestCreateUserRequest = entity.CreateUserRequest{
	FirstName: "Testy",
	LastName:  "McTestPants",
	Username:  "testy",
	Email:     "testy@example.com",
	Password:  "Password123",
}

var standardTestLoginRequest = entity.UserLoginRequest{
	Username: standardTestCreateUserRequest.Username,
	Password: standardTestCreateUserRequest.Password,
}

func TestCreateUser(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	{ // Success
		if _, err := s.CreateUser(standardTestCreateUserRequest); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		c.UserSettingsProviderValue.Provider.UpdateError = errors.New("foobar")
		if _, err := s.CreateUser(standardTestCreateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderValue.Provider.CreateError = errors.New("foobar")
		if _, err := s.CreateUser(standardTestCreateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := s.CreateUser(standardTestCreateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.CreateUser(standardTestCreateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}

func TestGetCurrentSession(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	userSession, err := s.SetupServer(standardTestCreateUserRequest)
	if err != nil {
		t.Fatalf("Failed to create user session: %s", err.Error())
	}

	{ // Success
		if _, err := s.GetCurrentSession(userSession.Token); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		if _, err := s.GetCurrentSession("INVALID_TOKEN"); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.SessionProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetCurrentSession(userSession.Token); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.SessionProviderError = errors.New("foobar")
		if _, err := s.GetCurrentSession(userSession.Token); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}

func TestGetUserContextByUUID(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	userSession, err := s.SetupServer(standardTestCreateUserRequest)
	if err != nil {
		t.Fatalf("Failed to create user session: %s", err.Error())
	}

	{ // Success
		if _, err := s.GetUserContextByUUID(userSession.Context.User.UUID); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		c.UserSettingsProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetUserContextByUUID(userSession.Context.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetUserContextByUUID(userSession.Context.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := s.GetUserContextByUUID(userSession.Context.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.GetUserContextByUUID(userSession.Context.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}

func TestLogin(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	_, err := s.SetupServer(standardTestCreateUserRequest)
	if err != nil {
		t.Fatalf("Failed to create user session: %s", err.Error())
	}

	{ // Success
		if _, err := s.Login(standardTestLoginRequest); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		c.SessionProviderValue.Provider.CreateError = errors.New("foobar")
		if _, err := s.Login(standardTestLoginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.SessionProviderError = errors.New("foobar")
		if _, err := s.Login(standardTestLoginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := s.Login(standardTestLoginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		if _, err := s.Login(entity.UserLoginRequest{
			Username: standardTestLoginRequest.Username,
			Password: "INVLAID_PASSWORD",
		}); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		if _, err := s.Login(entity.UserLoginRequest{
			Username: "INVALID_USERNAME",
			Password: "INVLAID_PASSWORD",
		}); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderValue.GetByUsernameError = errors.New("foobar")
		if _, err := s.Login(standardTestLoginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.Login(standardTestLoginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}
