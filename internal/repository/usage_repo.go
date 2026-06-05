package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
)

// UsageRepository handles usage record data access.
type UsageRepository struct {
	pool *pgxpool.Pool
}

// NewUsageRepository creates a new usage repository.
func NewUsageRepository(pool *pgxpool.Pool) *UsageRepository {
	return &UsageRepository{pool: pool}
}

// Create creates a new usage record.
func (r *UsageRepository) Create(ctx context.Context, record *model.UsageRecord) error {
	query := `
		INSERT INTO usage_records (id, organization_id, user_id, provider, model, input_tokens, output_tokens, cost, recorded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query, record.ID, record.OrganizationID, record.UserID, record.Provider, record.Model, record.InputTokens, record.OutputTokens, record.Cost, record.RecordedAt)
	return err
}

// GetByOrgIDAndPeriod retrieves usage records for an organization in a time period.
func (r *UsageRepository) GetByOrgIDAndPeriod(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]model.UsageRecord, error) {
	query := `
		SELECT id, organization_id, user_id, provider, model, input_tokens, output_tokens, cost, recorded_at
		FROM usage_records
		WHERE organization_id = $1 AND recorded_at >= $2 AND recorded_at <= $3
		ORDER BY recorded_at DESC
	`
	rows, err := r.pool.Query(ctx, query, orgID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []model.UsageRecord
	for rows.Next() {
		var rec model.UsageRecord
		if err := rows.Scan(&rec.ID, &rec.OrganizationID, &rec.UserID, &rec.Provider, &rec.Model, &rec.InputTokens, &rec.OutputTokens, &rec.Cost, &rec.RecordedAt); err != nil {
			return nil, err
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}

// GetUserUsageSummaries returns usage summaries for all users in an organization.
func (r *UsageRepository) GetUserUsageSummaries(ctx context.Context, orgID uuid.UUID, start, end time.Time) ([]model.UsageSummary, error) {
	query := `
		SELECT
			user_id,
			SUM(input_tokens) as total_input,
			SUM(output_tokens) as total_output,
			SUM(cost) as total_cost
		FROM usage_records
		WHERE organization_id = $1 AND recorded_at >= $2 AND recorded_at <= $3
		GROUP BY user_id
	`
	rows, err := r.pool.Query(ctx, query, orgID, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []model.UsageSummary
	for rows.Next() {
		var s model.UsageSummary
		if err := rows.Scan(&s.UserID, &s.TotalInput, &s.TotalOutput, &s.TotalCost); err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

// GetOrgTotalUsage returns total usage for an organization in a time period.
func (r *UsageRepository) GetOrgTotalUsage(ctx context.Context, orgID uuid.UUID, start, end time.Time) (int, int, float64, error) {
	query := `
		SELECT
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COALESCE(SUM(cost), 0)
		FROM usage_records
		WHERE organization_id = $1 AND recorded_at >= $2 AND recorded_at <= $3
	`
	var totalInput, totalOutput int
	var totalCost float64
	err := r.pool.QueryRow(ctx, query, orgID, start, end).Scan(&totalInput, &totalOutput, &totalCost)
	return totalInput, totalOutput, totalCost, err
}
