package hub

import (
	"os"
	"testing"
)

func TestNewProduction(t *testing.T) {
	targetKey := "LISTEN"
	targetVal := "0.0.0.0:1234"
	os.Setenv(targetKey, targetVal)
	service, _ := NewProduction()

	if service.Config.Listen != targetVal {
		t.Errorf(
			"%s has the value of: %s expected: %s",
			targetKey,
			service.Config.Listen,
			targetVal,
		)
	}
}
func TestNewProductionError(t *testing.T) {
	os.Setenv("DEV", "not a bool")
	_, err := NewProduction()

	if err == nil {
		t.Error("there should have been an error")
	}
}
