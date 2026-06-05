<<<<<<< HEAD
# AI Budget Splitter Backend

A Go-based backend service that acts as an API key rerouter, enabling organizations (up to 10 members) to track and split AI usage costs.

## Features

- **API Key Rerouting**: Users receive generated API keys to call the backend, which proxies requests to AI providers using real API keys
- **Organization Management**: Manage organizations with up to 10 members
- **Usage Tracking**: Track per-user AI usage (tokens, cost) with percentage breakdown
- **Billing Split**: Calculate proportional billing based on actual usage

## Tech Stack

- **Language**: Go 1.22+
- **Web Framework**: Gin
- **Database**: PostgreSQL with pgx
- **Configuration**: Viper

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 16+
- Make

### Setup

1. Clone the repository
2. Copy and configure environment variables:
   ```bash
   cp config.yaml config.local.yaml
   # Edit config.local.yaml with your settings
   ```

3. Set required environment variables:
   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ai_budget_splitter?sslmode=disable"
   export OPENAI_API_KEY="your-openai-key"
   export ANTHROPIC_API_KEY="your-anthropic-key"
   export JWT_SECRET="your-secret-key"
   ```

4. Run the application:
   ```bash
   make run
   ```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new organization
- `POST /api/v1/auth/login` - Login, returns API key
- `POST /api/v1/auth/keys` - Generate new API key
- `DELETE /api/v1/auth/keys/{key_id}` - Revoke API key
- `GET /api/v1/auth/keys` - List API keys

### Organization
- `GET /api/v1/org` - Get organization details
- `PUT /api/v1/org` - Update organization
- `POST /api/v1/org/members` - Add member (max 10)
- `DELETE /api/v1/org/members/{user_id}` - Remove member
- `GET /api/v1/org/members` - List members

### Proxy (AI Provider)
- `POST /api/v1/proxy/openai/*path` - Proxy to OpenAI
- `POST /api/v1/proxy/anthropic/*path` - Proxy to Anthropic

### Usage & Billing
- `GET /api/v1/usage` - Get current period usage
- `GET /api/v1/usage/history` - Get usage history
- `GET /api/v1/billing` - Get billing breakdown
- `GET /api/v1/billing/export` - Export billing data

## Development

### Make Commands

```bash
make build    # Build the server binary
make lint     # Run golangci-lint
make test     # Run tests with race detection
make run      # Run the server
make clean    # Remove build artifacts
make deps     # Download and tidy dependencies
```

### Running Tests

```bash
make test
```

## License

MIT
=======
# backend-ai-budget-splitter
>>>>>>> c63e8d279abf04b27a1cd0b3ba972ef248f85378
