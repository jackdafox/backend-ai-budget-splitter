package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString returns a SHA-256 hash of the input string.
func HashString(s string) string {
	hash := sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
