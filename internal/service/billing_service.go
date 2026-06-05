package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/repository"
)

// BillingService handles billing operations.
type BillingService struct {
	usageRepo *repository.UsageRepository
	userRepo  *repository.UserRepository
}

// NewBillingService creates a new billing service.
func NewBillingService(usageRepo *repository.UsageRepository, userRepo *repository.UserRepository) *BillingService {
	return &BillingService{
		usageRepo: usageRepo,
		userRepo:  userRepo,
	}
}

// BillingBreakdown represents the billing breakdown for an organization.
type BillingBreakdown struct {
	OrgID         uuid.UUID        `json:"org_id"`
	PeriodStart   time.Time        `json:"period_start"`
	PeriodEnd     time.Time        `json:"period_end"`
	TotalCost     float64          `json:"total_cost"`
	TotalInput    int              `json:"total_input_tokens"`
	TotalOutput   int              `json:"total_output_tokens"`
	MemberBillets []MemberBillable `json:"member_billets"`
}

// MemberBillable represents a member's billable amount.
type MemberBillable struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	Role         string    `json:"role"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	Cost         float64   `json:"cost"`
	Percentage   float64   `json:"percentage"`
	OwedAmount   float64   `json:"owed_amount"`
}

// GetBillingBreakdown returns the billing breakdown for a period.
func (s *BillingService) GetBillingBreakdown(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*BillingBreakdown, error) {
	// Get org usage totals
	totalInput, totalOutput, totalCost, err := s.usageRepo.GetOrgTotalUsage(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	// Get user summaries
	summaries, err := s.usageRepo.GetUserUsageSummaries(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	// Get users in org
	users, err := s.userRepo.GetByOrgID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	// Build user map
	userMap := make(map[uuid.UUID]*model.User)
	for i := range users {
		userMap[users[i].ID] = &users[i]
	}

	// Build billables
	var billables []MemberBillable
	for _, summary := range summaries {
		user, ok := userMap[summary.UserID]
		if !ok {
			continue
		}

		percentage := 0.0
		if totalCost > 0 {
			percentage = (summary.TotalCost / totalCost) * 100
		}

		billables = append(billables, MemberBillable{
			UserID:       summary.UserID,
			Email:        user.Email,
			Role:         string(user.Role),
			InputTokens:  summary.TotalInput,
			OutputTokens: summary.TotalOutput,
			Cost:         summary.TotalCost,
			Percentage:   percentage,
			OwedAmount:   summary.TotalCost, // Could apply markup/margin here
		})
	}

	return &BillingBreakdown{
		OrgID:         orgID,
		PeriodStart:   start,
		PeriodEnd:     end,
		TotalCost:     totalCost,
		TotalInput:    totalInput,
		TotalOutput:   totalOutput,
		MemberBillets: billables,
	}, nil
}

// ExportBillingData exports billing data as a structured format.
func (s *BillingService) ExportBillingData(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]byte, error) {
	breakdown, err := s.GetBillingBreakdown(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	// Could export as CSV, JSON, or PDF
	// For now, return JSON
	// Placeholder - actual JSON serialization would go here
	_ = breakdown
	return []byte{}, nil
}
