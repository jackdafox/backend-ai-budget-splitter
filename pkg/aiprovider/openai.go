package aiprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// OpenAIProvider implements the Provider interface for OpenAI.
type OpenAIProvider struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewOpenAIProvider creates a new OpenAI provider.
func NewOpenAIProvider(baseURL, apiKey string) *OpenAIProvider {
	return &OpenAIProvider{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// Name returns the provider name.
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Proxy proxies a request to OpenAI.
func (p *OpenAIProvider) Proxy(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*ProxyResponse, error) {
	url := p.BaseURL + "/" + path

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+p.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	usage := p.extractUsage(respBody)

	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	return &ProxyResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    respHeaders,
		Usage:      usage,
	}, nil
}

// extractUsage extracts usage from OpenAI response.
func (p *OpenAIProvider) extractUsage(body []byte) *UsageInfo {
	var resp struct {
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
		Model string `json:"model"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil
	}

	cost := float64(resp.Usage.PromptTokens)*0.00001 + float64(resp.Usage.CompletionTokens)*0.00003

	return &UsageInfo{
		InputTokens:  resp.Usage.PromptTokens,
		OutputTokens: resp.Usage.CompletionTokens,
		Cost:         cost,
		Model:        resp.Model,
	}
}
