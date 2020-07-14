package hub

import (
	"os"
	"testing"

	"github.com/fuzzingbits/hub/internal/container"
	"github.com/fuzzingbits/hub/internal/hubconfig"
)

func TestNewService(t *testing.T) {
	targetKey := "LISTEN"
	targetVal := "0.0.0.0:1234"
	os.Setenv(targetKey, targetVal)
	NewService(&hubconfig.Config{}, &container.Production{})
}
