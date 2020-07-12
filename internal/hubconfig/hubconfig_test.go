package hubconfig

import (
	"os"
	"testing"
)

func TestGetConfig(t *testing.T) {
	targetKey := "LISTEN"
	targetVal := "0.0.0.0:1234"
	os.Setenv(targetKey, targetVal)
	config, _ := GetConfig()

	if config.Listen != targetVal {
		t.Errorf(
			"%s has the value of: %s expected: %s",
			targetKey,
			config.Listen,
			targetVal,
		)
	}
}

func TestGetConfigError(t *testing.T) {
	os.Setenv("DEV", "not a bool")
	_, err := GetConfig()

	if err == nil {
		t.Error("there should have been an error")
	}
}
