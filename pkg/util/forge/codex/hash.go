package codex

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash a string with salt
func Hash(source string, salt string) string {
	hash := sha256.Sum256([]byte(source + salt))

	return hex.EncodeToString(hash[:])
}
