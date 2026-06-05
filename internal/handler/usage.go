package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/service"
)

// UsageHandler handles usage endpoints.
type UsageHandler struct {
	usageService   *service.UsageService
	billingService *service.BillingService
}

// NewUsageHandler creates a new usage handler.
func NewUsageHandler(usageService *service.UsageService, billingService *service.BillingService) *UsageHandler {
	return &UsageHandler{
		usageService:   usageService,
		billingService: billingService,
	}
}

// GetCurrentUsage handles GET /api/v1/usage.
func (h *UsageHandler) GetCurrentUsage(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	summary, err := h.usageService.GetCurrentPeriodUsage(c.Request.Context(), orgID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetUsageHistory handles GET /api/v1/usage/history.
func (h *UsageHandler) GetUsageHistory(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
			return
		}
	} else {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
			return
		}
	} else {
		end = time.Now()
	}

	summary, err := h.usageService.GetUsageHistory(c.Request.Context(), orgID.(uuid.UUID), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetBilling handles GET /api/v1/billing.
func (h *UsageHandler) GetBilling(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
			return
		}
	} else {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
			return
		}
	} else {
		end = time.Now()
	}

	breakdown, err := h.billingService.GetBillingBreakdown(c.Request.Context(), orgID.(uuid.UUID), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, breakdown)
}

// ExportBilling handles GET /api/v1/billing/export.
func (h *UsageHandler) ExportBilling(c *gin.Context) {
	orgID, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startStr := c.Query("start")
	endStr := c.Query("end")

	var start, end time.Time
	var err error

	if startStr != "" {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format"})
			return
		}
	} else {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	}

	if endStr != "" {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format"})
			return
		}
	} else {
		end = time.Now()
	}

	data, err := h.billingService.ExportBillingData(c.Request.Context(), orgID.(uuid.UUID), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
