# Configuration

The application is configured via `config.yaml` with environment variable interpolation.

## config.yaml Structure

```yaml
app:
  host: "0.0.0.0" # Listen address
  port: 8080          # Listen port
  env: "development"  # development or production

database:
  host: "localhost"
  port: 5432
  user: "postgres"
  password: "postgres"
  name: "ai_budget_splitter"
  sslmode: "disable"

ai_providers:
  openai:
    base_url: "https://api.openai.com/v1"
    api_key: "${OPENAI_API_KEY}"
  anthropic:
    base_url: "https://api.anthropic.com/v1"
    api_key: "${ANTHROPIC_API_KEY}"

auth:
  jwt_secret: "${JWT_SECRET}"
  key_expire_days: 90
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_URL` | PostgreSQL connection string (overrides database.*) | Yes |
| `OPENAI_API_KEY` | OpenAI API key for proxy requests | Yes |
| `ANTHROPIC_API_KEY` | Anthropic API key for proxy requests | Yes |
| `JWT_SECRET` | Secret key for signing JWT tokens | Yes |
| `APP_ENV` | Environment (development/production) | No |

## Example Setup

```bash
# .env file
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ai_budget_splitter?sslmode=disable"
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export JWT_SECRET="your-super-secret-key-at-least-32-chars"
```

## Production Considerations

1. **Use a real database** — PostgreSQL with SSL enabled in production
2. **Strong JWT secret** — Use a cryptographically random string (32+ chars)
3. **API keys** — Store in a secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.)
4. **TLS** — Run behind a reverse proxy (nginx, Caddy) with HTTPS
5. **Rate limiting** — Adjust rate limits in `internal/middleware/ratelimit.go`
