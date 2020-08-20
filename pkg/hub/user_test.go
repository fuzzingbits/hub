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
		invalidRealToken := "INVALID_REAL_TOKEN"
		c.SessionProviderValue.Set(invalidRealToken, "not a real user UUID")
		if _, err := s.GetCurrentSession(invalidRealToken); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.GetCurrentSession(userSession.Token); err == nil {
			t.Errorf("there should have been an error")
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

func TestUpdateUser(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	userSession, err := s.SetupServer(standardTestCreateUserRequest)
	if err != nil {
		t.Fatalf("Failed to create user session: %s", err.Error())
	}

	testUpdateUserRequest := entity.UpdateUserRequest{
		UUID: userSession.Context.User.UUID,
	}

	{ // Success
		if _, err := s.UpdateUser(testUpdateUserRequest); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		c.UserSettingsProviderValue.Provider.UpdateError = errors.New("foobar")
		if _, err := s.UpdateUser(testUpdateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderValue.Provider.UpdateError = errors.New("foobar")
		if _, err := s.UpdateUser(testUpdateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserSettingsProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.UpdateUser(testUpdateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderValue.Provider.GetByIDError = errors.New("foobar")
		if _, err := s.UpdateUser(testUpdateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserSettingsProviderError = errors.New("foobar")
		if _, err := s.UpdateUser(testUpdateUserRequest); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.UpdateUser(testUpdateUserRequest); err == nil {
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
		if _, err := s.GetUserContextByUUID("fake uuid"); err == nil {
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

func TestListUsers(t *testing.T) {
	c := container.NewMockable()
	s := NewService(&hubconfig.Config{}, c)

	_, _ = s.SetupServer(standardTestCreateUserRequest)
	_, _ = s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})

	{ // Success
		allUsers, err := s.ListUsers()
		if err != nil {
			t.Error(err)
		}

		if len(allUsers) != 2 {
			t.Errorf("Should bt 2 users, not %d", len(allUsers))
		}
	}

	{ // Error
		c.UserProviderValue.Provider.GetAllError = errors.New("foobar")
		if _, err := s.ListUsers(); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c.UserProviderError = errors.New("foobar")
		if _, err := s.ListUsers(); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}

func TestDeleteUser(t *testing.T) {
	testSetup := func() (*container.Mockable, *Service) {
		c := container.NewMockable()
		s := NewService(&hubconfig.Config{}, c)
		_, err := s.SetupServer(standardTestCreateUserRequest)
		if err != nil {
			t.Fatalf("Failed to create user session: %s", err.Error())
		}

		return c, s
	}

	{ // Success
		_, s := testSetup()
		userContext, _ := s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})
		if err := s.DeleteUser(userContext.User.UUID); err != nil {
			t.Error(err)
		}
	}

	{ // Error
		c, s := testSetup()
		userContext, _ := s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})
		c.UserSettingsProviderValue.Provider.DeleteError = errors.New("foobar")
		if err := s.DeleteUser(userContext.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c, s := testSetup()
		userContext, _ := s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})
		c.UserProviderValue.Provider.DeleteError = errors.New("foobar")
		if err := s.DeleteUser(userContext.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		_, s := testSetup()
		if err := s.DeleteUser("fake-uuid"); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c, s := testSetup()
		userContext, _ := s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})
		c.UserProviderValue.Provider.GetByIDError = errors.New("foobar")
		if err := s.DeleteUser(userContext.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c, s := testSetup()
		userContext, _ := s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})
		c.UserSettingsProviderError = errors.New("foobar")
		if err := s.DeleteUser(userContext.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}

	{ // Error
		c, s := testSetup()
		userContext, _ := s.CreateUser(entity.CreateUserRequest{Email: "foobar@example.com"})
		c.UserProviderError = errors.New("foobar")
		if err := s.DeleteUser(userContext.User.UUID); err == nil {
			t.Errorf("there should have been an error")
		}
	}
}
