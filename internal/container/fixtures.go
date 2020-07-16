package container

import (
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/provider/user"
	"github.com/fuzzingbits/hub/internal/provider/usersettings"
)

type dataProvider interface {
	AutoMigrate(clearExitstingDataAndCreateFixtures bool) error
}

func autoMigrateAll(providers []dataProvider, clearExitstingDataAndCreateFixtures bool) error {
	for _, provider := range providers {
		if err := provider.AutoMigrate(clearExitstingDataAndCreateFixtures); err != nil {
			return err
		}
	}

	return nil
}

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
			ThemeColor: "lime",
		})
	}
}
