package session

import (
	"time"
)

// Duration is the default session experation
const Duration = time.Hour * 24 * 6

// CookieName is the session cookie name
const CookieName = "HUB_SID"

// Provider is the Session Provider
type Provider interface {
	// Get a user UUID by token
	Get(token string) (string, error)
	// Set a session by token
	Set(token string, userUUI string) error
}
