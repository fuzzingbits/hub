package container

import (
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/provider/user"
	"github.com/fuzzingbits/hub/internal/provider/usersettings"
)

func (c *Production) createFixtures(
	userProvider user.Provider,
	userSettingProvider usersettings.Provider,
) {
	{ // Create Primary Test User
		uuid := "313efbe9-173b-4a1b-9a5b-7b69d95a66b9"
		userProvider.Create(entity.User{
			UUID:      uuid,
			FirstName: "Aaron",
			LastName:  "Ellington",
		})

		userSettingProvider.Save(uuid, entity.UserSettings{
			ThemeColor: "#00BFFF",
		})
	}
}
