package aiprovider

import (
	"context"
	"io"
)

// Provider defines the interface for AI providers.
type Provider interface {
	// Name returns the provider name.
	Name() string

	// Proxy proxies a request to the provider.
	Proxy(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*ProxyResponse, error)
}

// ProxyResponse represents a response from an AI provider.
type ProxyResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	Usage      *UsageInfo
}

// UsageInfo contains usage information.
type UsageInfo struct {
	InputTokens  int
	OutputTokens int
	Cost         float64
	Model        string
}
