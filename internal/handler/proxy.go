package handler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/service"
)

// ProxyHandler handles AI provider proxy endpoints.
type ProxyHandler struct {
	proxyService *service.ProxyService
}

// NewProxyHandler creates a new proxy handler.
func NewProxyHandler(proxyService *service.ProxyService) *ProxyHandler {
	return &ProxyHandler{proxyService: proxyService}
}

// ProxyOpenAI handles POST /api/v1/proxy/openai/*path.
func (h *ProxyHandler) ProxyOpenAI(c *gin.Context) {
	path := c.Param("path")
	orgID, _ := c.Get("org_id")
	userID, _ := c.Get("user_id")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	req := &service.ProxyRequest{
		Method:  c.Request.Method,
		Path:    path,
		OrgID:   orgID.(uuid.UUID),
		UserID:  userID.(uuid.UUID),
		Body:    bytes.NewReader(body),
		Headers: extractHeaders(c),
	}

	resp, err := h.proxyService.ProxyOpenAI(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for k, v := range resp.Headers {
		c.Header(k, v)
	}
	c.Data(resp.StatusCode, "application/json", resp.Body)
}

// ProxyAnthropic handles POST /api/v1/proxy/anthropic/*path.
func (h *ProxyHandler) ProxyAnthropic(c *gin.Context) {
	path := c.Param("path")
	orgID, _ := c.Get("org_id")
	userID, _ := c.Get("user_id")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	req := &service.ProxyRequest{
		Method:  c.Request.Method,
		Path:    path,
		OrgID:   orgID.(uuid.UUID),
		UserID:  userID.(uuid.UUID),
		Body:    bytes.NewReader(body),
		Headers: extractHeaders(c),
	}

	resp, err := h.proxyService.ProxyAnthropic(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for k, v := range resp.Headers {
		c.Header(k, v)
	}
	c.Data(resp.StatusCode, "application/json", resp.Body)
}

// extractHeaders extracts relevant headers from the request.
func extractHeaders(c *gin.Context) map[string]string {
	headers := make(map[string]string)
	skipHeaders := map[string]bool{
		"content-type":    true,
		"content-length":  true,
		"authorization":   true,
		"x-api-key":       true,
		"anthropic-version": true,
	}

	for k, v := range c.Request.Header {
		if !skipHeaders[k] && len(v) > 0 {
			headers[k] = v[0]
		}
	}
	return headers
}
