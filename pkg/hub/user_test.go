package hub

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
	"github.com/fuzzingbits/hub/pkg/provider/session"
)

func TestGetCurrentSession(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{RollbarToken: "fake-token"}, c)

	// Create Fixture User
	targetSession, _ := s.SetupServer(
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

	{ // No Cookie
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}

	req.AddCookie(&http.Cookie{
		Name:    session.CookieName,
		Value:   targetSession.Token,
		Expires: time.Now().Add(time.Minute),
	})

	{ // Test Success!
		userSession, err := s.GetCurrentSession(req)
		if err != nil {
			t.Error(err)
			return
		}

		if !reflect.DeepEqual(userSession, targetSession) {
			t.Errorf(
				"[Session Did Not Match] returned: %s expected: %s",
				userSession,
				targetSession,
			)
		}
	}

	{
		c.SessionProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetCurrentSession(req); err == nil {
			t.Error("there should have been an error")
		}
	}

	{
		c.SessionProviderError = errors.New("foobar")
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

	c.UserProviderValue.Provider.CreateError = errors.New("foobar")
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

func TestLoginMissingRequestBody(t *testing.T) {
	// c := container.NewMockable()
	// s := NewService(&hubconfig.Config{}, c)

	// if _, err := loginTestHelper(t, s, nil); err == nil {
	// 	t.Errorf("there should have been an error")
	// }
}

func TestLoginBadRequestBody(t *testing.T) {
	// loginRequest := entity.UserLoginRequest{
	// 	Username: "testy",
	// 	Password: "Password123",
	// }

	// loginRequestBytes, _ := json.Marshal(loginRequest)

	// c := container.NewMockable()
	// s := NewService(&hubconfig.Config{}, c)

	// if _, err := loginTestHelper(t, s, loginRequestBytes[:1]); err == nil {
	// 	t.Errorf("there should have been an error")
	// }

}

func TestLoginSuccess(t *testing.T) {
	loginRequest := entity.UserLoginRequest{
		Username: "testy",
		Password: "Password123",
	}

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	if _, err := s.Login(loginRequest); err != nil {
		t.Error(err)
	}
}

func TestLoginIncorrectPassword(t *testing.T) {
	loginRequest := entity.UserLoginRequest{
		Username: "testy",
		Password: "NotTheCorrectPassword",
	}

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	if _, err := s.Login(loginRequest); err == nil {
		t.Errorf("there should have been an error")
	}
}

func TestLoginUserNotFound(t *testing.T) {
	loginRequest := entity.UserLoginRequest{
		Username: "testy",
		Password: "NotTheCorrectPassword",
	}

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	if _, err := s.Login(loginRequest); err == nil {
		t.Errorf("there should have been an error")
	}
}

func TestGetUserContextByUUID(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	createdUser, _ := s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	if _, err := s.GetUserContextByUUID(createdUser.User.UUID); err != nil {
		t.Error(err)
	}

	c.UserSettingsProviderValue.Provider.GetByIDError = errors.New("foobar")
	if _, err := s.GetUserContextByUUID(createdUser.User.UUID); err == nil {
		t.Errorf("there should have been an error")
	}

	c.UserProviderValue.Provider.GetByIDError = errors.New("foobar")
	if _, err := s.GetUserContextByUUID(createdUser.User.UUID); err == nil {
		t.Errorf("there should have been an error")
	}

	c.UserProviderError = errors.New("foobar")
	if _, err := s.GetUserContextByUUID(createdUser.User.UUID); err == nil {
		t.Errorf("there should have been an error")
	}
}

func TestLoginErrors(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	loginRequest := entity.UserLoginRequest{
		Username: "testy",
		Password: "Password123",
	}

	s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	{
		c.SessionProviderValue.Provider.CreateError = errors.New("foobar")
		if _, err := s.Login(loginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.SessionProviderError = errors.New("foobar")
		if _, err := s.Login(loginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := s.Login(loginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.UserProviderValue.GetByUsernameError = errors.New("foobar")
		if _, err := s.Login(loginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.UserProviderError = errors.New("foobar")
		if _, err := s.Login(loginRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}
