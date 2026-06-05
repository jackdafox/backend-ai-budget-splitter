package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/model"
	"github.com/jackdafox/backend-ai-budget-splitter/internal/repository"
)

// ProxyService handles AI provider proxy operations.
type ProxyService struct {
	usageRepo  *repository.UsageRepository
	openAIURL  string
	openAIKey  string
	anthropicURL string
	anthropicKey string
}

// NewProxyService creates a new proxy service.
func NewProxyService(usageRepo *repository.UsageRepository, openAIURL, openAIKey, anthropicURL, anthropicKey string) *ProxyService {
	return&ProxyService{
		usageRepo:     usageRepo,
		openAIURL:     openAIURL,
		openAIKey:     openAIKey,
		anthropicURL: anthropicURL,
		anthropicKey: anthropicKey,
	}
}

// ProxyRequest represents a proxy request to an AI provider.
type ProxyRequest struct {
	Method         string
	Path           string
	OrgID          uuid.UUID
	UserID         uuid.UUID
	Body           io.Reader
	Headers        map[string]string
}

// ProxyResponse represents a response from an AI provider.
type ProxyResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Usage      *UsageInfo
}

// UsageInfo contains usage information extracted from AI provider response.
type UsageInfo struct {
	InputTokens  int
	OutputTokens int
	Cost         float64
	Model        string
}

// ProxyOpenAI proxies a request to OpenAI.
func (s *ProxyService) ProxyOpenAI(ctx context.Context, req *ProxyRequest) (*ProxyResponse, error) {
	url := s.openAIURL + "/" + req.Path

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.openAIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	usage := s.extractOpenAIUsage(body, resp.Header)

	// Record usage
	if usage != nil {
		record := &model.UsageRecord{
			ID:             uuid.New(),
			OrganizationID: req.OrgID,
			UserID:         req.UserID,
			Provider:       "openai",
			Model:          usage.Model,
			InputTokens:    usage.InputTokens,
			OutputTokens:   usage.OutputTokens,
			Cost:           usage.Cost,
			RecordedAt:     time.Now(),
		}
		_ = s.usageRepo.Create(ctx, record)
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &ProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    headers,
		Usage:      usage,
	}, nil
}

// ProxyAnthropic proxies a request to Anthropic.
func (s *ProxyService) ProxyAnthropic(ctx context.Context, req *ProxyRequest) (*ProxyResponse, error) {
	url := s.anthropicURL + "/" + req.Path

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("x-api-key", s.anthropicKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	usage := s.extractAnthropicUsage(body, resp.Header)

	// Record usage
	if usage != nil {
		record := &model.UsageRecord{
			ID:             uuid.New(),
			OrganizationID: req.OrgID,
			UserID:         req.UserID,
			Provider:       "anthropic",
			Model:          usage.Model,
			InputTokens:    usage.InputTokens,
			OutputTokens:   usage.OutputTokens,
			Cost:           usage.Cost,
			RecordedAt:     time.Now(),
		}
		_ = s.usageRepo.Create(ctx, record)
	}

	headers := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	return &ProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       body,
		Headers:    headers,
		Usage:      usage,
	}, nil
}

// extractOpenAIUsage extracts usage information from OpenAI response.
func (s *ProxyService) extractOpenAIUsage(body []byte, header http.Header) *UsageInfo {
	var resp struct {
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens       int `json:"total_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil
	}

	// Calculate cost (approximate OpenAI pricing)
	cost := float64(resp.Usage.PromptTokens)*0.00001 + float64(resp.Usage.CompletionTokens)*0.00003

	return &UsageInfo{
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		Cost:         cost,
		Model:        resp.Model,
	}
}

// extractAnthropicUsage extracts usage information from Anthropic response.
func (s *ProxyService) extractAnthropicUsage(body []byte, header http.Header) *UsageInfo {
	var resp struct {
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil
	}

	// Calculate cost (approximate Anthropic pricing)
	cost := float64(resp.Usage.InputTokens)*0.000015 + float64(resp.Usage.OutputTokens)*0.000075

	return &UsageInfo{
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		Cost:         cost,
		Model:        resp.Model,
	}
}

// ProxyStreamOpenAI handles streaming responses from OpenAI.
func (s *ProxyService) ProxyStreamOpenAI(ctx context.Context, req *ProxyRequest) (*http.Response, error) {
	url := s.openAIURL + "/" + req.Path

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+s.openAIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	return client.Do(httpReq)
}

// ProxyStreamAnthropic handles streaming responses from Anthropic.
func (s *ProxyService) ProxyStreamAnthropic(ctx context.Context, req *ProxyRequest) (*http.Response, error) {
	url := s.anthropicURL + "/" + req.Path

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("x-api-key", s.anthropicKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Accept", "text/event-stream")
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 120 * time.Second}
	return client.Do(httpReq)
}

// CopyBody reads and resets the body buffer.
func CopyBody(body io.Reader) (io.ReadCloser, *bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	reader := io.TeeReader(body, buf)
	return io.NopCloser(reader), buf, nil
}
