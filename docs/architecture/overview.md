# Architecture Overview

## System Design

The AI Budget Splitter backend follows a **layered architecture** with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Layer (Gin)                        │
│   Routes → Middleware → Handlers                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Service Layer                            │
│   AuthService │ ProxyService │ UsageService │ BillingService │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Repository Layer                          │
│   OrgRepo │ UserRepo │ APIKeyRepo │ UsageRepo               │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                  Database (PostgreSQL) │
└─────────────────────────────────────────────────────────────┘
```

## Request Flow

### Proxy Request Flow

```
Client (with generated API key)
    │
    ▼
┌─────────────────┐
│  Auth Middleware │ ← Validates JWT
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ Rate Limiter     │ ← Per-user rate limits
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ Proxy Handler    │ ← Routes to ProxyService
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ Proxy Service    │ ← Forwards to AI provider
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ Usage Repository│ ← Records usage after response
└─────────────────┘
    │
    ▼
┌─────────────────┐
│ AI Provider     │ ← OpenAI / Anthropic
│ (real API key)   │
└─────────────────┘
```

### Authentication Flow

```
1. Register → Creates Org + Admin User → Returns JWT
2. Login → Validates credentials → Returns JWT
3. Generate API Key → Returns one-time API key
4. API calls use JWT for org-level operations
5. Proxy calls use JWT for user attribution
```

## Key Components

### Handlers (`internal/handler/`)
Handle HTTP requests/responses. Each handler group maps to a resource:
- `auth.go` — Registration, login, API key management
- `org.go` — Organization and member management
- `proxy.go` — AI provider proxy endpoints
- `usage.go` — Usage queries and billing

### Services (`internal/service/`)
Contain business logic:
- `AuthService` — JWT generation, API key hashing, credential validation
- `ProxyService` — Request forwarding, usage extraction from responses
- `UsageService` — Usage aggregation and percentage calculation
- `BillingService` — Cost breakdown per member

### Repositories (`internal/repository/`)
Data access layer using pgx:
- `OrgRepository` — Organization CRUD
- `UserRepository` — User CRUD, org member lookup
- `APIKeyRepository` — API key storage (bcrypt hashed), validation
- `UsageRepository` — Usage record storage and aggregation queries

### Middleware (`internal/middleware/`)
- `AuthMiddleware` — JWT validation, sets `user_id`, `org_id` in context
- `RateLimiter` — Token bucket rate limiting per user
- `Logger` — Request logging

## Database Schema

```
organizations
├── id (UUID, PK)
├── name
├── created_at
└── updated_at

users
├── id (UUID, PK)
├── organization_id (FK → organizations)
├── email (UNIQUE)
├── password_hash
├── role (admin/member)
└── created_at

api_keys
├── id (UUID, PK)
├── user_id (FK → users)
├── organization_id (FK → organizations)
├── key_hash (bcrypt)
├── key_prefix (first 8 chars)
├── name
├── last_used_at
├── expires_at
└── created_at

usage_records
├── id (UUID, PK)
├── organization_id (FK → organizations)
├── user_id (FK → users)
├── provider (openai/anthropic)
├── model
├── input_tokens
├── output_tokens
├── cost
└── recorded_at
```

## Security Model

1. **Passwords** — bcrypt hashed before storage
2. **API Keys** — bcrypt hashed; only prefix stored in responses
3. **JWT** — HS256 signed; contains user ID, org ID, role
4. **Rate Limiting** — 100 requests/minute per user (configurable)
5. **Input Validation** — Gin binding with structural validation

## Configuration

Managed via Viper with `config.yaml` + environment variable overrides. See [Configuration](../getting-started/configuration.md).
