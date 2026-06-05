package model

import (
	"time"

	"github.com/google/uuid"
)

// MemberRole represents the role of a member in an organization.
type MemberRole string

const (
	RoleAdmin  MemberRole = "admin"
	RoleMember MemberRole = "member"
)

// User represents a user in the system.
type User struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"org_id"`
	Email          string    `json:"email"`
	PasswordHash   string    `json:"-"`
	Role           MemberRole `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}
