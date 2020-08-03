package hub

import (
	"testing"

	"github.com/fuzzingbits/hub/pkg/container"
	"github.com/fuzzingbits/hub/pkg/hubconfig"
)

func TestNewService(t *testing.T) {
	c := container.NewMockable()
	NewService(&hubconfig.Config{RollbarToken: "FAKE_TOKEN"}, c)
}
