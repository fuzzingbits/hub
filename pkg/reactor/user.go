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
		Username:  dbUser.Username,
		Email:     dbUser.Email,
	}
}

// CreateUserRequestToDBUser does with it says
func CreateUserRequestToDBUser(request entity.CreateUserRequest) entity.DatabaseUser {
	newUUID := uuid.New().String()
	password := codex.Hash(request.Password, newUUID)
	return entity.DatabaseUser{
		UUID:      newUUID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Username:  request.Username,
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
