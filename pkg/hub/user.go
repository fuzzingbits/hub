package hub

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/provider/user"
	"github.com/fuzzingbits/hub/pkg/reactor"
	"github.com/fuzzingbits/hub/pkg/util/forge/codex"
	"github.com/google/uuid"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// ErrInvalidLogin is when the login credentials are incorrect
var ErrInvalidLogin = errors.New("Invalid Login")

// CreateUser creates a User
func (s *Service) CreateUser(request entity.CreateUserRequest) (entity.UserSession, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	dbUser := reactor.CreateUserRequestToDBUser(request)

	if err := userProvider.Create(&dbUser); err != nil {
		return entity.UserSession{}, err
	}

	// Setup new UserSettings with defaults
	userSettings := entity.UserSettings{
		ThemeColor: "tomato",
	}

	if err := userSettingsProvider.Save(dbUser.UUID, userSettings); err != nil {
		return entity.UserSession{}, err
	}

	return entity.UserSession{
		User:     reactor.DatabaseUserToEntity(dbUser),
		Settings: userSettings,
	}, nil
}

// GetCurrentSession gets the current session
func (s *Service) GetCurrentSession(r *http.Request) (entity.UserSession, error) {
	return s.GetUserSessionByUUID(r.Header.Get("UUID"))
}

// GetUserSessionByUUID by UUID and get the full session
func (s *Service) GetUserSessionByUUID(uuid string) (entity.UserSession, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	// TODO: create actual session management and stop doing this terrible thing
	user, err := userProvider.GetByUUID(uuid)
	if err != nil {
		return entity.UserSession{}, err
	}

	userSettings, err := userSettingsProvider.GetByUUID(user.UUID)
	if err != nil {
		return entity.UserSession{}, err
	}

	return entity.UserSession{
		User:     reactor.DatabaseUserToEntity(user),
		Settings: userSettings,
	}, nil
}

// Login attempts to create a session
func (s *Service) Login(w http.ResponseWriter, req *http.Request) (entity.UserSession, error) {
	if req.Body == nil {
		return entity.UserSession{}, fmt.Errorf("Missing Request Body")
	}

	// Build the login request
	var loginRequest entity.UserLoginRequest
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		return entity.UserSession{}, err
	}

	// Get the UserProvider
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	// Query the user that we are trying to login as
	potentialUser, err := userProvider.GetByUsername(loginRequest.Username)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return entity.UserSession{}, ErrInvalidLogin
		}
		return entity.UserSession{}, err
	}

	// Verify the password
	providedPassword := codex.Hash(loginRequest.Password, potentialUser.UUID)
	if providedPassword != potentialUser.Password {
		return entity.UserSession{}, ErrInvalidLogin
	}

	// Get the full session based on the verified user UUID
	userSession, err := s.GetUserSessionByUUID(potentialUser.UUID)
	if err != nil {
		return entity.UserSession{}, err
	}

	// Get the session provider
	sessionProvider, err := s.container.SessionProvider()
	if err != nil {
		return entity.UserSession{}, err
	}

	// Build a new Session token
	newSessionToken := uuid.New().String()

	// Save the session
	if err := sessionProvider.Set(newSessionToken, userSession); err != nil {
		return entity.UserSession{}, err
	}

	// Set the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:    session.CookieName,
		Value:   newSessionToken,
		Expires: time.Now().Add(session.Duration),
	})

	return userSession, nil
}
