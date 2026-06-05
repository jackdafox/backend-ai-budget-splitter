package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/repository"
)

// OrgHandler handles organization endpoints.
type OrgHandler struct {
	orgRepo  *repository.OrganizationRepository
	userRepo *repository.UserRepository
}

// NewOrgHandler creates a new org handler.
func NewOrgHandler(orgRepo *repository.OrganizationRepository, userRepo *repository.UserRepository) *OrgHandler {
	return &OrgHandler{
		orgRepo:  orgRepo,
		userRepo: userRepo,
	}
}

// GetOrg handles GET /api/v1/org.
func (h *OrgHandler) GetOrg(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	org, err := h.orgRepo.GetByID(c.Request.Context(), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrg handles PUT /api/v1/org.
func (h *OrgHandler) UpdateOrg(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := h.orgRepo.GetByID(c.Request.Context(), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	org.Name = req.Name
	if err := h.orgRepo.Update(c.Request.Context(), org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// AddMember handles POST /api/v1/org/members.
func (h *OrgHandler) AddMember(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check member limit
	count, err := h.userRepo.CountByOrgID(c.Request.Context(), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count >= 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "organization member limit reached (max 10)"})
		return
	}

	// Check if email exists
	if _, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	}

	// Create user
	user := &model.User{
		ID:             uuid.New(),
		OrganizationID: orgID.(uuid.UUID),
		Email:          req.Email,
		PasswordHash:   req.Password, // Should be hashed in real impl
		Role:           model.RoleMember,
		CreatedAt:      time.Now(),
	}
	if req.Role == "admin" {
		user.Role = model.RoleAdmin
	}

	if err := h.userRepo.Create(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// RemoveMember handles DELETE /api/v1/org/members/:user_id.
func (h *OrgHandler) RemoveMember(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	if err := h.userRepo.Delete(c.Request.Context(), userID); err != nil {
		if err == repository.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "member removed"})
}

// ListMembers handles GET /api/v1/org/members.
func (h *OrgHandler) ListMembers(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := h.userRepo.GetByOrgID(c.Request.Context(), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"members": users})
}
