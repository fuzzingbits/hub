package session

import (
	"time"

	"github.com/fuzzingbits/hub/pkg/entity"
)

// Duration is the default session experation
const Duration = time.Hour * 24 * 6

// CookieName is the session cookie name
const CookieName = "HUB_SID"

// Provider is the Session Provider
type Provider interface {
	// Get a session by token
	Get(token string) (entity.Session, error)
	// Set a session by token
	Set(token string, session entity.Session) error
}
