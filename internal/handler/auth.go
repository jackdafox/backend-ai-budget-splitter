package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/service"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GenerateAPIKey handles POST /api/v1/auth/keys.
func (h *AuthHandler) GenerateAPIKey(c *gin.Context) {
	var req service.GenerateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	apiKey, err := h.authService.GenerateAPIKey(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"api_key":     apiKey.KeyHash, // Only returned once
		"key_prefix":  apiKey.KeyPrefix,
		"key_id":      apiKey.ID,
		"name":        apiKey.Name,
		"expires_at":  apiKey.ExpiresAt,
		"created_at":  apiKey.CreatedAt,
	})
}

// RevokeAPIKey handles DELETE /api/v1/auth/keys/:key_id.
func (h *AuthHandler) RevokeAPIKey(c *gin.Context) {
	keyID, err := uuid.Parse(c.Param("key_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key_id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.authService.RevokeAPIKey(c.Request.Context(), keyID, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "API key revoked"})
}

// ListAPIKeys handles GET /api/v1/auth/keys.
func (h *AuthHandler) ListAPIKeys(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	keys, err := h.authService.ListAPIKeys(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"api_keys": keys})
}
