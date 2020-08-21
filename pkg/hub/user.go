package hub

import (
	"errors"

	"github.com/fuzzingbits/hub/pkg/provider/session"
	"github.com/fuzzingbits/hub/pkg/provider/user"
	"github.com/fuzzingbits/hub/pkg/reactor"
	"github.com/fuzzingbits/hub/pkg/util/forge/codex"
	"github.com/fuzzingbits/hub/pkg/util/forge/gol"
	"github.com/google/uuid"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// DefaultThemeColor for new users
var DefaultThemeColor = "#00bfff"

// ErrInvalidLogin is when the login credentials are incorrect
var ErrInvalidLogin = errors.New("Invalid Login")

// ErrMissingValidSession is when there is no valid session
var ErrMissingValidSession = errors.New("No Valid Session")

// ErrRecordNotFound is when no record is found
var ErrRecordNotFound = errors.New("Record not found")

// ListUsers gets all users
func (s *Service) ListUsers() ([]entity.User, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return nil, err
	}

	dbUsers, err := userProvider.GetAll()
	if err != nil {
		return nil, err
	}

	entityUsers := []entity.User{}
	for _, dbUser := range dbUsers {
		entityUsers = append(entityUsers, reactor.DatabaseUserToEntity(dbUser))
	}

	return entityUsers, nil
}

// UpdateUser updates the full user context
func (s *Service) UpdateUser(request entity.UpdateUserRequest) (entity.UserContext, error) {
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return entity.UserContext{}, err
	}

	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return entity.UserContext{}, err
	}

	dbUser, err := userProvider.GetByUUID(request.UUID)
	if err != nil {
		return entity.UserContext{}, err
	}

	userSettings, err := userSettingsProvider.GetByUUID(request.UUID)
	if err != nil {
		return entity.UserContext{}, err
	}

	// Apply the update data to the existing data
	reactor.ApplyUserUpdateRequest(request, &dbUser, &userSettings)

	if err := userProvider.Update(&dbUser); err != nil {
		return entity.UserContext{}, err
	}

	if err := userSettingsProvider.Save(request.UUID, userSettings); err != nil {
		return entity.UserContext{}, err
	}

	userContext, err := s.GetUserContextByUUID(request.UUID)
	gol.PanicOnError(err)

	return userContext, nil
}

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
		ThemeColor: DefaultThemeColor,
	}

	if err := userSettingsProvider.Save(dbUser.UUID, userSettings); err != nil {
		return entity.UserContext{}, err
	}

	return entity.UserContext{
		User:     reactor.DatabaseUserToEntity(dbUser),
		Settings: userSettings,
	}, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(uuid string) error {
	// Get the user provider
	userProvider, err := s.container.UserProvider()
	if err != nil {
		return err
	}

	// Get the user settings provider
	userSettingsProvider, err := s.container.UserSettingsProvider()
	if err != nil {
		return err
	}

	// Get the exiting user by the provided UUID
	dbUser, err := userProvider.GetByUUID(uuid)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return ErrRecordNotFound
		}

		return err
	}

	// Delete the user
	if err := userProvider.Delete(dbUser); err != nil {
		return err
	}

	// Delete the user settings
	if err := userSettingsProvider.Delete(uuid); err != nil {
		return err
	}

	return nil
}

// GetCurrentSession gets the current session
func (s *Service) GetCurrentSession(token string) (entity.Session, error) {
	sessionProvider, err := s.container.SessionProvider()
	if err != nil {
		return entity.Session{}, err
	}

	userUUID, err := sessionProvider.Get(token)
	if err != nil {
		if errors.Is(err, session.ErrNotFound) {
			return entity.Session{}, ErrMissingValidSession
		}

		return entity.Session{}, err
	}

	userContext, err := s.GetUserContextByUUID(userUUID)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return entity.Session{}, ErrMissingValidSession
		}

		return entity.Session{}, err
	}

	return entity.Session{
		Token:   token,
		Context: userContext,
	}, nil
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

	databaseUser, err := userProvider.GetByUUID(uuid)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return entity.UserContext{}, ErrRecordNotFound
		}

		return entity.UserContext{}, err
	}

	userSettings, err := userSettingsProvider.GetByUUID(databaseUser.UUID)
	if err != nil {
		return entity.UserContext{}, err
	}

	return entity.UserContext{
		User:     reactor.DatabaseUserToEntity(databaseUser),
		Settings: userSettings,
	}, nil
}

// SaveSession saves a session
func (s *Service) SaveSession(sessionID string, userContext entity.UserContext) (entity.Session, error) {
	// Get the session provider
	sessionProvider, err := s.container.SessionProvider()
	if err != nil {
		return entity.Session{}, err
	}

	// Build a new Session token
	userSession := entity.Session{
		Token:   sessionID,
		Context: userContext,
	}

	// Save the session
	if err := sessionProvider.Set(userSession.Token, userContext.User.UUID); err != nil {
		return entity.Session{}, err
	}

	return userSession, nil
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

	return s.SaveSession(uuid.New().String(), userContext)
}
