package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/repository"
)

// UsageService handles usage tracking operations.
type UsageService struct {
	usageRepo *repository.UsageRepository
	userRepo  *repository.UserRepository
}

// NewUsageService creates a new usage service.
func NewUsageService(usageRepo *repository.UsageRepository, userRepo *repository.UserRepository) *UsageService {
	return &UsageService{
		usageRepo: usageRepo,
		userRepo:  userRepo,
	}
}

// GetCurrentPeriodUsage returns usage for the current billing period.
func (s *UsageService) GetCurrentPeriodUsage(ctx context.Context, orgID uuid.UUID) (*model.OrgUsageSummary, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Second)

	return s.getUsageSummary(ctx, orgID, start, end)
}

// GetUsageHistory returns usage for a specific time period.
func (s *UsageService) GetUsageHistory(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*model.OrgUsageSummary, error) {
	return s.getUsageSummary(ctx, orgID, start, end)
}

// getUsageSummary calculates usage summary for an organization.
func (s *UsageService) getUsageSummary(ctx context.Context, orgID uuid.UUID, start, end time.Time) (*model.OrgUsageSummary, error) {
	// Get total org usage
	totalInput, totalOutput, totalCost, err := s.usageRepo.GetOrgTotalUsage(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	// Get per-user summaries
	summaries, err := s.usageRepo.GetUserUsageSummaries(ctx, orgID, start, end)
	if err != nil {
		return nil, err
	}

	// Calculate percentages
	for i := range summaries {
		if totalCost > 0 {
			summaries[i].UsagePercent = (summaries[i].TotalCost / totalCost) * 100
		}
	}

	return &model.OrgUsageSummary{
		OrganizationID: orgID,
		TotalCost:      totalCost,
		TotalInput:     totalInput,
		TotalOutput:    totalOutput,
		UserSummaries:  summaries,
	}, nil
}

// GetUserUsage returns usage for a specific user.
func (s *UsageService) GetUserUsage(ctx context.Context, userID uuid.UUID, start, end time.Time) (*model.UsageSummary, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	summaries, err := s.usageRepo.GetUserUsageSummaries(ctx, user.OrganizationID, start, end)
	if err != nil {
		return nil, err
	}

	for _, s := range summaries {
		if s.UserID == userID {
			return &s, nil
		}
	}

	return &model.UsageSummary{
		UserID:       userID,
		TotalInput:   0,
		TotalOutput:  0,
		TotalCost:    0,
		UsagePercent: 0,
	}, nil
}
