package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
)

// APIKeyRepository handles API key data access.
type APIKeyRepository struct {
	pool *pgxpool.Pool
}

// NewAPIKeyRepository creates a new API key repository.
func NewAPIKeyRepository(pool *pgxpool.Pool) *APIKeyRepository {
	return &APIKeyRepository{pool: pool}
}

// Create creates a new API key.
func (r *APIKeyRepository) Create(ctx context.Context, key *model.APIKey) error {
	query := `
		INSERT INTO api_keys (id, user_id, organization_id, key_hash, key_prefix, name, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query, key.ID, key.UserID, key.OrganizationID, key.KeyHash, key.KeyPrefix, key.Name, key.ExpiresAt, key.CreatedAt)
	return err
}

// GetByID retrieves an API key by ID.
func (r *APIKeyRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.APIKey, error) {
	query := `
		SELECT id, user_id, organization_id, key_hash, key_prefix, name, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE id = $1
	`
	var key model.APIKey
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&key.ID, &key.UserID, &key.OrganizationID, &key.KeyHash, &key.KeyPrefix, &key.Name, &key.LastUsedAt, &key.ExpiresAt, &key.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &key, err
}

// GetByHash retrieves an API key by its hash.
func (r *APIKeyRepository) GetByHash(ctx context.Context, hash string) (*model.APIKey, error) {
	query := `
		SELECT id, user_id, organization_id, key_hash, key_prefix, name, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE key_hash = $1
	`
	var key model.APIKey
	err := r.pool.QueryRow(ctx, query, hash).Scan(
		&key.ID, &key.UserID, &key.OrganizationID, &key.KeyHash, &key.KeyPrefix, &key.Name, &key.LastUsedAt, &key.ExpiresAt, &key.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &key, err
}

// GetByUserID retrieves all API keys for a user.
func (r *APIKeyRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
	query := `
		SELECT id, user_id, organization_id, key_hash, key_prefix, name, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []model.APIKey
	for rows.Next() {
		var key model.APIKey
		if err := rows.Scan(&key.ID, &key.UserID, &key.OrganizationID, &key.KeyHash, &key.KeyPrefix, &key.Name, &key.LastUsedAt, &key.ExpiresAt, &key.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

// GetByOrgID retrieves all API keys for an organization.
func (r *APIKeyRepository) GetByOrgID(ctx context.Context, orgID uuid.UUID) ([]model.APIKey, error) {
	query := `
		SELECT id, user_id, organization_id, key_hash, key_prefix, name, last_used_at, expires_at, created_at
		FROM api_keys
		WHERE organization_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.pool.Query(ctx, query, orgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []model.APIKey
	for rows.Next() {
		var key model.APIKey
		if err := rows.Scan(&key.ID, &key.UserID, &key.OrganizationID, &key.KeyHash, &key.KeyPrefix, &key.Name, &key.LastUsedAt, &key.ExpiresAt, &key.CreatedAt); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, rows.Err()
}

// UpdateLastUsed updates the last_used_at timestamp.
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE api_keys SET last_used_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

// Delete deletes an API key.
func (r *APIKeyRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_keys WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
