package reactor

import (
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/forge/codex"
	"github.com/fuzzingbits/hub/internal/provider/user"
	"github.com/google/uuid"
)

// DatabaseUserToEntity does what is says
func DatabaseUserToEntity(dbUser user.User) entity.User {
	return entity.User{
		UUID:      dbUser.UUID,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
	}
}

func CreateUserRequestToDBUser(request entity.CreateUserRequest) user.User {
	newUUID := uuid.New().String()
	password := codex.Hash(request.Password, newUUID)
	return user.User{
		UUID:      newUUID,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Username:  request.Username,
		Password:  password,
	}
}

// EntityToDatabaseUser does what is says
func EntityToDatabaseUser(entityUser entity.User) user.User {
	return user.User{
		UUID:      entityUser.UUID,
		FirstName: entityUser.FirstName,
		LastName:  entityUser.LastName,
	}
}
