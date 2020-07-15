package hub

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/entity"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestGetCurrentSession(t *testing.T) {
	container := container.NewMockable()
	service := NewService(&hubconfig.Config{}, container)

	targetUser := entity.User{
		UUID:      "313efbe9-173b-4a1b-9a5b-7b69d95a66b9",
		FirstName: "Aaron",
		LastName:  "Ellington",
	}

	container.UserProviderValue.Create(targetUser)

	session, err := service.GetCurrentSession(&http.Request{})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(session, entity.UserSession{
		User: targetUser,
	}) {
		t.Errorf("session did not match")

	}
}

func TestGetCurrentSessionContainerError(t *testing.T) {
	container := container.NewMockable()
	service := NewService(&hubconfig.Config{}, container)

	container.UserSettingsProviderValue = nil
	if _, err := service.GetCurrentSession(&http.Request{}); err == nil {
		t.Errorf("should have returned an error")
	}

	container.UserProviderValue = nil
	if _, err := service.GetCurrentSession(&http.Request{}); err == nil {
		t.Errorf("should have returned an error")
	}
}

func TestGetCurrentSessionProviderError(t *testing.T) {
	container := container.NewMockable()
	service := NewService(&hubconfig.Config{}, container)

	if _, err := service.GetCurrentSession(&http.Request{}); err == nil {
		t.Errorf("should have returned an error")
	}

	targetUser := entity.User{
		UUID:      "313efbe9-173b-4a1b-9a5b-7b69d95a66b9",
		FirstName: "Aaron",
		LastName:  "Ellington",
	}
	container.UserProviderValue.Create(targetUser)

	container.UserSettingsProviderValue.GetByUUIDError = errors.New("random error")
	if _, err := service.GetCurrentSession(&http.Request{}); err == nil {
		t.Errorf("should have returned an error")
	}
}
