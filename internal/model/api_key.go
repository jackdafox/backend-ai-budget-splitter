package model

import (
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key for a user.
type APIKey struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	OrganizationID uuid.UUID  `json:"org_id"`
	KeyHash        string     `json:"-"`
	KeyPrefix      string     `json:"key_prefix"` // First 8 chars for display
	Name           string     `json:"name"`
	LastUsedAt     *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt      *time.Time `json:"expires_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}
