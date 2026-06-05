package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// OrganizationRepository handles organization data access.
type OrganizationRepository struct {
	pool *pgxpool.Pool
}

// NewOrganizationRepository creates a new organization repository.
func NewOrganizationRepository(pool *pgxpool.Pool) *OrganizationRepository {
	return &OrganizationRepository{pool: pool}
}

// Create creates a new organization.
func (r *OrganizationRepository) Create(ctx context.Context, org *model.Organization) error {
	query := `
		INSERT INTO organizations (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pool.Exec(ctx, query, org.ID, org.Name, org.CreatedAt, org.UpdatedAt)
	return err
}

// GetByID retrieves an organization by ID.
func (r *OrganizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM organizations
		WHERE id = $1
	`
	var org model.Organization
	err := r.pool.QueryRow(ctx, query, id).Scan(&org.ID, &org.Name, &org.CreatedAt, &org.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &org, err
}

// Update updates an organization.
func (r *OrganizationRepository) Update(ctx context.Context, org *model.Organization) error {
	query := `
		UPDATE organizations
		SET name = $2, updated_at = $3
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, org.ID, org.Name, org.UpdatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete deletes an organization.
func (r *OrganizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM organizations WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
