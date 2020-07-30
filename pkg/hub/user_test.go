package hub

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
)

func TestGetCurrentSession(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{RollbarToken: "fake-token"}, c)

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
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	if _, err := loginTestHelper(t, s, nil); err == nil {
		t.Errorf("there should have been an error")
	}
}

func TestLoginBadRequestBody(t *testing.T) {
	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "Password123",
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	if _, err := loginTestHelper(t, s, loginRequestBytes[:1]); err == nil {
		t.Errorf("there should have been an error")
	}

}

func TestLoginSuccess(t *testing.T) {
	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "Password123",
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	if _, err := loginTestHelper(t, s, loginRequestBytes); err != nil {
		t.Error(err)
	}
}

func TestLoginIncorrectPassword(t *testing.T) {
	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "NotTheCorrectPassword",
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
		t.Errorf("there should have been an error")
	}
}

func TestLoginUserNotFound(t *testing.T) {
	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "NotTheCorrectPassword",
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
		t.Errorf("there should have been an error")
	}
}

func TestLoginErrors(t *testing.T) {
	loginRequest := entity.CreateUserRequest{
		Username: "testy",
		Password: "Password123",
	}

	loginRequestBytes, _ := json.Marshal(loginRequest)

	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	s.CreateUser(entity.CreateUserRequest{
		FirstName: "Testy",
		LastName:  "McTestPants",
		Username:  "testy",
		Email:     "testy@example.com",
		Password:  "Password123",
	})

	{
		c.SessionProviderValue.Provider.CreateError = errors.New("foobar")
		if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.SessionProviderError = errors.New("foobar")
		if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.UserProviderValue.GetByUsernameError = errors.New("foobar")
		if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{
		c.UserProviderError = errors.New("foobar")
		if _, err := loginTestHelper(t, s, loginRequestBytes); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}

func loginTestHelper(t *testing.T, s *Service, body []byte) (entity.UserSession, error) {
	var payload io.Reader
	if body != nil {
		payload = bytes.NewReader(body)
	}

	x := httptest.NewRecorder()
	r, _ := http.NewRequest("", "/", payload)

	return s.Login(x, r)
}
