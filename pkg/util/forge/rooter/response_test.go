package rooter

import (
	"testing"
)

func TestResponses(t *testing.T) {
	ResponseInternalServerError()
	ResponseMethodNotAllowed()
	ResponseNotFound()
}