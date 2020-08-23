package reactor

import (
	"github.com/fuzzingbits/hub/pkg/entity"
	"github.com/fuzzingbits/hub/pkg/util/forge/codex"
	"github.com/google/uuid"
)

// DatabaseUserToEntity does with it says
func DatabaseUserToEntity(dbUser entity.DatabaseUser) entity.User {
	return entity.User{
		UUID:      dbUser.UUID,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Email:     dbUser.Email,
	}
}

// UserCreateRequestToDBUser does with it says
func UserCreateRequestToDBUser(request entity.UserCreateRequest) entity.DatabaseUser {
	newUUID := uuid.New().String()
	password := codex.Hash(request.Password, newUUID)
	return entity.DatabaseUser{
		UUID:      newUUID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  password,
	}
}

// EntityToDatabaseUser does with it says
func EntityToDatabaseUser(entityUser entity.User) entity.DatabaseUser {
	return entity.DatabaseUser{
		UUID:      entityUser.UUID,
		FirstName: entityUser.FirstName,
		LastName:  entityUser.LastName,
	}
}

// ApplyUserUpdateRequest applies a request on to a user context
func ApplyUserUpdateRequest(request entity.UserUpdateRequest, dbUser *entity.DatabaseUser, userSettings *entity.UserSettings) {
	dbUser.FirstName = request.FirstName
	dbUser.LastName = request.LastName
	dbUser.Email = request.Email

	userSettings.ThemeColorDark = request.ThemeColorDark
	userSettings.ThemeColorLight = request.ThemeColorLight
}
