package model

import (
	"time"

	"github.com/google/uuid"
)

// UsageRecord represents an AI usage record.
type UsageRecord struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"org_id"`
	UserID         uuid.UUID `json:"user_id"`
	Provider       string    `json:"provider"` // "openai", "anthropic"
	Model          string    `json:"model"`
	InputTokens    int       `json:"input_tokens"`
	OutputTokens   int       `json:"output_tokens"`
	Cost           float64   `json:"cost"`
	RecordedAt     time.Time `json:"recorded_at"`
}

// UsageSummary represents a usage summary for a user.
type UsageSummary struct {
	UserID         uuid.UUID `json:"user_id"`
	TotalInput     int       `json:"total_input_tokens"`
	TotalOutput    int       `json:"total_output_tokens"`
	TotalCost      float64   `json:"total_cost"`
	UsagePercent   float64   `json:"usage_percent"` // Percentage of org total
}

// OrgUsageSummary represents usage summary for an entire organization.
type OrgUsageSummary struct {
	OrganizationID uuid.UUID      `json:"org_id"`
	TotalCost      float64        `json:"total_cost"`
	TotalInput     int            `json:"total_input_tokens"`
	TotalOutput    int            `json:"total_output_tokens"`
	UserSummaries  []UsageSummary `json:"user_summaries"`
}
