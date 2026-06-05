# Proxy API

Proxy requests to AI providers (OpenAI, Anthropic) using the organization's shared API keys. Usage is tracked per user.

## Proxy to OpenAI

**Endpoint:** `POST /api/v1/proxy/openai/*path`

```bash
curl -X POST http://localhost:8080/api/v1/proxy/openai/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_HERE" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}],
    "temperature": 0.7
  }'
```

**Response** — Forwards OpenAI response directly:

```json
{
  "id": "chatcmpl-...",
  "object": "chat.completion",
  "model": "gpt-3.5-turbo-0301",
  "choices": [...],
  "usage": {
    "prompt_tokens": 10,
    "completion_tokens": 20,
    "total_tokens": 30
  }
}
```

## Proxy to Anthropic

**Endpoint:** `POST /api/v1/proxy/anthropic/*path`

```bash
curl -X POST http://localhost:8080/api/v1/proxy/anthropic/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_HERE" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-sonnet-20240229",
    "messages": [{"role": "user", "content": "Hello!"}],
    "max_tokens": 256
  }'
```

**Response** — Forwards Anthropic response directly:

```json
{
  "id": "msg_...",
  "type": "message",
  "model": "claude-3-sonnet-20240229",
  "content": [...],
  "usage": {
    "input_tokens": 10,
    "output_tokens": 20
  }
}
```

## Supported Paths

### OpenAI
- `/chat/completions` — Chat completions
- `/completions` — Text completions
- `/embeddings` — Embeddings
- `/images/generations` — Image generation

### Anthropic
- `/messages` — Messages API (Claude)
- `/messages-stream` — Streaming messages

## Usage Tracking

Usage is automatically recorded after each proxy request:

| Field | Description |
|-------|-------------|
| `provider` | `openai` or `anthropic` |
| `model` | Model name from response |
| `input_tokens` | Tokens in the request |
| `output_tokens` | Tokens in the response |
| `cost` | Calculated cost (USD, approximate) |

## Rate Limits

- **100 requests per minute** per user (configurable in `ratelimit.go`)
- Exceeded limits return `429 Too Many Requests`

## Error Responses

| Status | Error | Cause |
|--------|-------|-------|
| `400` | `failed to read body` | Invalid request body |
| `429` | `rate limit exceeded` | Too many requests |
| `500` | `...` | AI provider error (forwarded) |
