package usersettings

import (
	"context"

	"github.com/fuzzingbits/hub/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type databaseUserSettings struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty"`
	UUID         string              `bson:"uuid"`
	UserSettings entity.UserSettings `bson:"userSettings"`
}

// DatabaseProvider is a usersettings.Provider that uses a database
type DatabaseProvider struct {
	Collection *mongo.Collection
}

// AutoMigrate the data connection
func (d *DatabaseProvider) AutoMigrate(clearExitstingData bool) error {
	if clearExitstingData {
		return d.Collection.Drop(context.TODO())
	}

	return nil
}

// GetByUUID gets a UserSettings by User.UUID
func (d *DatabaseProvider) GetByUUID(uuid string) (entity.UserSettings, error) {
	// Create the filter
	filterCursor, err := d.Collection.Find(context.TODO(), bson.M{"uuid": uuid})
	if err != nil {
		return entity.UserSettings{}, err
	}

	// Get all results
	var resultsTarget []databaseUserSettings
	if err := filterCursor.All(context.TODO(), &resultsTarget); err != nil {
		return entity.UserSettings{}, err
	}

	// Return just our single result
	for _, e := range resultsTarget {
		return e.UserSettings, nil
	}

	return entity.UserSettings{}, nil
}

// Save a UserSettings
func (d *DatabaseProvider) Save(uuid string, userSettings entity.UserSettings) error {
	d.Collection.InsertOne(context.TODO(), databaseUserSettings{
		UUID:         uuid,
		UserSettings: userSettings,
	})

	return nil
}

// Delete a UserSettings by UUID
func (d *DatabaseProvider) Delete(uuid string) error {
	return nil
}
