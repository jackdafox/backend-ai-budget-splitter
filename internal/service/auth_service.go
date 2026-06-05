package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/repository"
)

// AuthService handles authentication operations.
type AuthService struct {
	orgRepo       *repository.OrganizationRepository
	userRepo      *repository.UserRepository
	apiKeyRepo    *repository.APIKeyRepository
	jwtSecret     []byte
	keyExpireDays int
}

// NewAuthService creates a new auth service.
func NewAuthService(orgRepo *repository.OrganizationRepository, userRepo *repository.UserRepository, apiKeyRepo *repository.APIKeyRepository, jwtSecret string, keyExpireDays int) *AuthService {
	return &AuthService{
		orgRepo:       orgRepo,
		userRepo:      userRepo,
		apiKeyRepo:    apiKeyRepo,
		jwtSecret:     []byte(jwtSecret),
		keyExpireDays: keyExpireDays,
	}
}

// RegisterRequest represents a registration request.
type RegisterRequest struct {
	OrgName   string `json:"org_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// RegisterResponse represents a registration response.
type RegisterResponse struct {
	OrgID  uuid.UUID `json:"org_id"`
	UserID uuid.UUID `json:"user_id"`
	Token  string   `json:"token"`
}

// Register creates a new organization and admin user.
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check if email already exists
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("email already registered")
	}

	// Create organization
	org := &model.Organization{
		ID:        uuid.New(),
		Name:      req.OrgName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.orgRepo.Create(ctx, org); err != nil {
		return nil, err
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create admin user
	user := &model.User{
		ID:             uuid.New(),
		OrganizationID: org.ID,
		Email:          req.Email,
		PasswordHash:   string(hash),
		Role:           model.RoleAdmin,
		CreatedAt:      time.Now(),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		OrgID:  org.ID,
		UserID: user.ID,
		Token:  token,
	}, nil
}

// LoginRequest represents a login request.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a login response.
type LoginResponse struct {
	Token string       `json:"token"`
	User  *model.User  `json:"user"`
}

// Login authenticates a user and returns a JWT.
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// GenerateAPIKeyRequest represents an API key generation request.
type GenerateAPIKeyRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

// GenerateAPIKey generates a new API key for a user.
func (s *AuthService) GenerateAPIKey(ctx context.Context, req *GenerateAPIKeyRequest) (*model.APIKey, error) {
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// Generate random key
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, err
	}
	key := hex.EncodeToString(keyBytes)

	// Hash the key for storage
	hash, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Calculate expiry
	var expiresAt *time.Time
	if s.keyExpireDays > 0 {
		t := time.Now().AddDate(0, 0, s.keyExpireDays)
		expiresAt = &t
	}

	apiKey := &model.APIKey{
		ID:             uuid.New(),
		UserID:         user.ID,
		OrganizationID: user.OrganizationID,
		KeyHash:        string(hash),
		KeyPrefix:      key[:8],
		Name:           req.Name,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	// Return key with plaintext (only time this is returned)
	apiKey.KeyHash = key
	return apiKey, nil
}

// ValidateAPIKey validates an API key and returns the associated key record.
func (s *AuthService) ValidateAPIKey(ctx context.Context, key string) (*model.APIKey, error) {
	// We need to check all keys since we only store hashes
	// In production, you'd want to use a faster lookup with a key prefix index
	keys, err := s.apiKeyRepo.GetByOrgID(ctx, uuid.Nil)
	if err != nil && err != repository.ErrNotFound {
		return nil, err
	}

	// For simplicity, iterate through orgs - in production use a direct hash lookup
	// This is a limitation to address: store key_prefix separately for fast lookup
	_ = keys
	return nil, errors.New("api key validation requires key prefix lookup")
}

// RevokeAPIKey revokes an API key.
func (s *AuthService) RevokeAPIKey(ctx context.Context, keyID, userID uuid.UUID) error {
	key, err := s.apiKeyRepo.GetByID(ctx, keyID)
	if err != nil {
		return err
	}
	if key.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.apiKeyRepo.Delete(ctx, keyID)
}

// ListAPIKeys lists all API keys for a user.
func (s *AuthService) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
	return s.apiKeyRepo.GetByUserID(ctx, userID)
}

// generateJWT generates a JWT for a user.
func (s *AuthService) generateJWT(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"org":   user.OrganizationID.String(),
		"email": user.Email,
		"role":  string(user.Role),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
