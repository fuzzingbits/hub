package hub

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/provider/user"
	"github.com/fuzzingbits/hub/pkg/reactor"
	"github.com/fuzzingbits/hub/pkg/util/forge/codex"
	"github.com/google/uuid"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// ErrInvalidLogin is when the login credentials are incorrect
var ErrInvalidLogin = errors.New("Invalid Login")

// ErrMissingValidSession is when there is no valid session
var ErrMissingValidSession = errors.New("No Valid Session")

// CreateUser creates a User
func (s *Service) CreateUser(request entity.CreateUserRequest) (entity.UserContext, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserContext{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserContext{}, err
	}

	dbUser := reactor.CreateUserRequestToDBUser(request)

	if err := userProvider.Create(&dbUser); err != nil {
		return entity.UserContext{}, err
	}

	// Setup new UserSettings with defaults
	userSettings := entity.UserSettings{
		ThemeColor: "tomato",
	}

	if err := userSettingsProvider.Save(dbUser.UUID, userSettings); err != nil {
		return entity.UserContext{}, err
	}

	return entity.UserContext{
		User:     reactor.DatabaseUserToEntity(dbUser),
		Settings: userSettings,
	}, nil
}

// GetCurrentSession gets the current session
func (s *Service) GetCurrentSession(token string) (entity.Session, error) {
	sessionProvider, err := s.container.SessionProvider()
	if err != nil {
		return entity.Session{}, err
	}

	userSession, err := sessionProvider.Get(token)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return entity.Session{}, ErrMissingValidSession
		}

		return entity.Session{}, err
	}

	return userSession, nil
}

// GetUserContextByUUID by UUID and get the full conext
func (s *Service) GetUserContextByUUID(uuid string) (entity.UserContext, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserContext{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserContext{}, err
	}

	user, err := userProvider.GetByUUID(uuid)
	if err != nil {
		return entity.UserContext{}, err
	}

	userSettings, err := userSettingsProvider.GetByUUID(user.UUID)
	if err != nil {
		return entity.UserContext{}, err
	}

	return entity.UserContext{
		User:     reactor.DatabaseUserToEntity(user),
		Settings: userSettings,
	}, nil
}

// Login attempts to create a session
func (s *Service) Login(loginRequest entity.UserLoginRequest) (entity.Session, error) {
	// Get the UserProvider
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.Session{}, err
	}

	// Query the user that we are trying to login as
	potentialUser, err := userProvider.GetByUsername(loginRequest.Username)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return entity.Session{}, ErrInvalidLogin
		}

		return entity.Session{}, err
	}

	// Verify the password
	providedPassword := codex.Hash(loginRequest.Password, potentialUser.UUID)
	if providedPassword != potentialUser.Password {
		return entity.Session{}, ErrInvalidLogin
	}

	// Get the full context based on the verified user UUID
	userContext, err := s.GetUserContextByUUID(potentialUser.UUID)
	if err != nil {
		return entity.Session{}, err
	}

	// Get the session provider
	sessionProvider, err := s.container.SessionProvider()
	if err != nil {
		return entity.Session{}, err
	}

	// Build a new Session token
	userSession := entity.Session{
		Token:   uuid.New().String(),
		Context: userContext,
	}

	// Save the session
	if err := sessionProvider.Set(userSession.Token, userSession); err != nil {
		return entity.Session{}, err
	}

	return userSession, nil
}
