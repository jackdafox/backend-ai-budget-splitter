can# AI Budget Splitter Backend - Project Plan

## Overview

A Go-based backend service that acts as an API key rerouter, enabling organizations (up to 10 members) to track and split AI usage costs. Users receive generated API keys to call the backend, which proxies requests to AI providers using real API keys, tracking per-user usage percentages.

---

## 1. Project Structure

```
backend-ai-budget-splitter/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading
│   ├── handler/
│   │   ├── auth.go # API key generation/validation
│   │   ├── proxy.go             # AI provider proxy handlers
│   │   └── usage.go             # Usage tracking endpoints
│   ├── middleware/
│   │   ├── auth.go # API key authentication middleware
│   │   ├── ratelimit.go         # Rate limiting
│   │   └── logging.go           # Request logging
│   ├── model/
│   │   ├── user.go              # User model
│   │   ├── api_key.go           # API key model
│   │   ├── organization.go       # Organization model
│   │   └── usage.go             # Usage record model
│   ├── repository/
│   │   ├── user_repo.go         # User data access
│   │   ├── api_key_repo.go      # API key data access
│   │   ├── org_repo.go          # Organization data access
│   │   └── usage_repo.go        # Usage data access
│   ├── service/
│   │   ├── auth_service.go      # API key generation/validation
│   │   ├── proxy_service.go     # AI provider proxy logic
│   │   ├── usage_service.go      # Usage tracking& calculation
│   │   └── billing_service.go    # Bill splitting logic
│   └── database/
│       ├── postgres.go          # PostgreSQL connection
│       └── migrations/          # SQL migrations
├── pkg/
│   ├── aiprovider/
│   │   ├── openai.go            # OpenAI provider adapter
│   │   ├── anthropic.go        # Anthropic provider adapter
│   │   └── provider.go          # Provider interface
│   └── utils/
│       ├── hash.go              # Hashing utilities
│       └── response.go          # HTTP response helpers
├── api/
│   └── openapi.yaml             # OpenAPI specification
├── .github/
│   └── workflows/
│       ├── ci.yml               # CI pipeline (lint + test)
│       └── cd.yml               # Deploy pipeline
├── .golangci.yml               # Go linter configuration
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 2. Core Features

### 2.1 API Key Management
- [ ] Generate unique API keys per user/organization
- [ ] Validate API keys on each request
- [ ] Support key revocation and rotation
- [ ] Store hashed keys (never plaintext)

### 2.2 Organization Management
- [ ] Create organization (max 10 members)
- [ ] Add/remove members
- [ ] List organization members
- [ ] Set member roles (admin/member)

### 2.3 AI Provider Proxy
- [ ] Proxy requests to OpenAI API
- [ ] Proxy requests to Anthropic API
- [ ] Support multiple provider abstraction
- [ ] Forward original request headers/body
- [ ] Handle streaming responses

### 2.4 Usage Tracking
- [ ] Track tokens used per user per request
- [ ] Store usage records with timestamps
- [ ] Calculate usage percentages per member
- [ ] Generate usage reports

### 2.5 Billing Split
- [ ] Calculate proportional split based on usage %
- [ ] Generate invoice data per member
- [ ] Export billing summaries

---

## 3. API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register new organization |
| POST | `/api/v1/auth/login` | Login, returns API key |
| POST | `/api/v1/auth/keys` | Generate new API key |
| DELETE | `/api/v1/auth/keys/{key_id}` | Revoke API key |
| GET | `/api/v1/auth/keys` | List API keys |

### Organization
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/org` | Get organization details |
| PUT | `/api/v1/org` | Update organization |
| POST | `/api/v1/org/members` | Add member (max 10) |
| DELETE | `/api/v1/org/members/{user_id}` | Remove member |
| GET | `/api/v1/org/members` | List members |

### Proxy (AI Provider)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/proxy/openai/*path` | Proxy to OpenAI |
| POST | `/api/v1/proxy/anthropic/*path` | Proxy to Anthropic |

### Usage & Billing
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/usage` | Get current period usage |
| GET | `/api/v1/usage/history` | Get usage history |
| GET | `/api/v1/billing` | Get billing breakdown |
| GET | `/api/v1/billing/export` | Export billing data |

---

## 4. Data Models

### Organization
```go
type Organization struct {
    ID        uuid.UUID  `json:"id"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

### User
```go
type User struct {
    ID             uuid.UUID      `json:"id"`
    OrganizationID uuid.UUID      `json:"org_id"`
    Email          string         `json:"email"`
    PasswordHash   string         `json:"-"`
    Role           MemberRole     `json:"role"`
    CreatedAt      time.Time      `json:"created_at"`
}

type MemberRole string
const (
    RoleAdmin  MemberRole = "admin"
    RoleMember MemberRole = "member"
)
```

### APIKey
```go
type APIKey struct {
    ID             uuid.UUID  `json:"id"`
    UserID         uuid.UUID  `json:"user_id"`
    OrganizationID uuid.UUID  `json:"org_id"`
    KeyHash        string     `json:"-"`
    KeyPrefix      string     `json:"key_prefix"` // First 8 chars for display
    Name           string     `json:"name"`
    LastUsedAt     *time.Time `json:"last_used_at"`
    ExpiresAt      *time.Time `json:"expires_at"`
    CreatedAt      time.Time  `json:"created_at"`
}
```

### UsageRecord
```go
type UsageRecord struct {
    ID             uuid.UUID `json:"id"`
    OrganizationID uuid.UUID `json:"org_id"`
    UserID         uuid.UUID `json:"user_id"`
    Provider       string    `json:"provider"` // "openai", "anthropic"
    Model          string    `json:"model"`
    InputTokens    int       `json:"input_tokens"`
    OutputTokens   int       `json:"output_tokens"`
    Cost           float64   `json:"cost"`
    RecordedAt     time.Time `json:"recorded_at"`
}
```

---

## 5. Configuration

### config.yaml
```yaml
app:
  host: "0.0.0.0"
  port: 8080
  env: "development"

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

---

## 6. GitHub Actions CI Pipeline

### Workflow: `ci.yml`
```yaml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run golangci-lint
        run: make lint

  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: test_db
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Run tests
        run: make test
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable

  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build
        run: make build
```

---

## 7. Makefile Commands

```makefile
.PHONY: build lint test test-coverage run clean migrate

build:
	go build -o bin/server ./cmd/server

lint:
	golangci-lint run ./...

test:
	go test -v -race ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

run:
	go run ./cmd/server

clean:
	rm -rf bin/

migrate:
	go run ./cmd/migrate

deps:
	go mod download
	go mod tidy
```

---

## 8. Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `OPENAI_API_KEY` | OpenAI API key for proxy | Yes |
| `ANTHROPIC_API_KEY` | Anthropic API key for proxy | Yes |
| `JWT_SECRET` | Secret for JWT signing | Yes |
| `APP_ENV` | Environment (development/production) | No |

---

## 9. Security Considerations

- [ ] API keys stored as bcrypt hashes only
- [ ] Rate limiting per API key (configurable)
- [ ] Request logging (excluding sensitive data)
- [ ] CORS configuration for allowed origins
- [ ] Input validation on all endpoints
- [ ] SQL injection prevention via parameterized queries
- [ ] HTTPS enforcement in production

---

## 10. Implementation Phases

### Phase 1: Foundation
- [ ] Project setup with Go modules
- [ ] Configuration loading
- [ ] Database connection & migrations
- [ ] Basic HTTP server setup

### Phase 2: Authentication
- [ ] User registration/login
- [ ] API key generation & validation
- [ ] JWT middleware

### Phase 3: Proxy Core
- [ ] OpenAI proxy endpoint
- [ ] Anthropic proxy endpoint
- [ ] Usage extraction from responses
- [ ] Streaming response handling

### Phase 4: Usage Tracking
- [ ] Usage record storage
- [ ] Percentage calculation
- [ ] Usage query endpoints

### Phase 5: Billing
- [ ] Bill splitting calculation
- [ ] Export functionality
- [ ] Organization member management

### Phase 6: Polish
- [ ] Rate limiting
- [ ] Comprehensive tests
- [ ] Documentation
- [ ] CI/CD setup

---

## 11. Testing Strategy

### Unit Tests
- Service layer business logic
- Utility functions
- Hashing/encryption

### Integration Tests
- Database repositories
- API endpoints
- Proxy functionality

### Test Coverage Target
- Minimum 70% coverage
- Critical paths: 90%+

---

## 12. Dependencies

```go
// Core
github.com/golang-jwt/jwt/v5
github.com/google/uuid

// Web Framework
github.com/gin-gonic/gin

// Database
github.com/jackc/pgx/v5
github.com/golang-migrate/migrate/v4

// Configuration
github.com/spf13/viper

// Security
golang.org/x/crypto/bcrypt

// HTTP Client
github.com/go-resty/resty/v2

// Testing
github.com/stretchr/testify
github.com/vektra/mockery/v2
```

---

## 13. TODO Checklist

- [ ] Initialize Go module
- [ ] Create project directory structure
- [ ] Set up configuration with Viper
- [ ] Implement database connection
- [ ] Write database migrations
- [ ] Implement user/auth models & repository
- [ ] Implement API key service
- [ ] Set up Gin router with middleware
- [ ] Implement auth handlers
- [ ] Implement proxy handlers
- [ ] Implement usage tracking
- [ ] Implement billing service
- [ ] Write unit tests
- [ ] Write integration tests
- [ ] Configure golangci-lint
- [ ] Set up GitHub Actions workflow
- [ ] Create Makefile
- [ ] Write README documentation
