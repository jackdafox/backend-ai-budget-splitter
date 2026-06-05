# Quick Start

Get the backend running and make your first API call in under 5 minutes.

## 1. Start the Server

```bash
make run
```

Server runs on `http://localhost:8080`.

## 2. Register an Organization

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "org_name": "My Team",
    "email": "admin@example.com",
    "password": "securepassword123"
  }'
```

Response:

```json
{
  "org_id": "uuid-here",
  "user_id": "uuid-here",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Save the `token` — you'll need it for subsequent requests.

## 3. Generate an API Key

```bash
curl -X POST http://localhost:8080/api/v1/auth/keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "user_id": "YOUR_USER_ID",
    "name": "My API Key"
  }'
```

Response (save the `api_key` — it's shown only once):

```json
{
  "api_key": "abc123...",
  "key_prefix": "abc12345",
  "key_id": "uuid-here",
  "name": "My API Key",
  "expires_at": "2024-09-01T00:00:00Z",
  "created_at": "2024-06-01T00:00:00Z"
}
```

## 4. Make a Proxy Request

```bash
curl -X POST http://localhost:8080/api/v1/proxy/openai/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

The request is proxied to OpenAI, and usage is tracked for your organization.

## 5. Check Usage

```bash
curl http://localhost:8080/api/v1/usage \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

##6. View Billing Breakdown

```bash
curl http://localhost:8080/api/v1/billing \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Next Steps

- [Add organization members](/api/organization.md)
- [Understand the architecture](/architecture/overview.md)
- [Explore all API endpoints](/api/auth.md)
