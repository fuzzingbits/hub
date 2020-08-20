package session

import (
	"errors"
	"time"
)

// Duration is the default session experation
const Duration = time.Hour * 24 * 6

// CookieName is the session cookie name
const CookieName = "HUB_SID"

// ErrNotFound is when the session can not be found
var ErrNotFound = errors.New("Not Found")

// Provider is the Session Provider
type Provider interface {
	// Get a user UUID by token
	Get(token string) (string, error)
	// Set a session by token
	Set(token string, userUUI string) error
}
