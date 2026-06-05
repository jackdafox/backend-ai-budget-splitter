package aiprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// AnthropicProvider implements the Provider interface for Anthropic.
type AnthropicProvider struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewAnthropicProvider creates a new Anthropic provider.
func NewAnthropicProvider(baseURL, apiKey string) *AnthropicProvider {
	return &AnthropicProvider{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// Name returns the provider name.
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Proxy proxies a request to Anthropic.
func (p *AnthropicProvider) Proxy(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*ProxyResponse, error) {
	url := p.BaseURL + "/" + path

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("x-api-key", p.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("anthropic-version", "2023-06-01")
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

// extractUsage extracts usage from Anthropic response.
func (p *AnthropicProvider) extractUsage(body []byte) *UsageInfo {
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

	cost := float64(resp.Usage.InputTokens)*0.000015 + float64(resp.Usage.OutputTokens)*0.000075

	return &UsageInfo{
		InputTokens:  resp.Usage.InputTokens,
		OutputTokens: resp.Usage.OutputTokens,
		Cost:         cost,
		Model:        resp.Model,
	}
}
