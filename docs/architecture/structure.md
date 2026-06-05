# Project Structure

```
backend-ai-budget-splitter/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration loading (Viper)
│   ├── database/
│   │   ├── postgres.go          # PostgreSQL connection pool
│   │   └── migrations/          # SQL migration files
│   ├── handler/
│   │   ├── auth.go             # Auth endpoints
│   │   ├── org.go              # Organization endpoints
│   │   ├── proxy.go            # AI provider proxy
│   │   └── usage.go            # Usage & billing endpoints
│   ├── middleware/
│   │   ├── auth.go             # JWT authentication
│   │   ├── ratelimit.go        # Rate limiting
│   │   └── logging.go          # Request logging
│   ├── model/
│   │   ├── organization.go     # Organization model
│   │   ├── user.go            # User model & roles
│   │   ├── api_key.go         # API key model
│   │   └── usage.go           # Usage record model
│   ├── repository/
│   │   ├── org_repo.go        # Organization data access
│   │   ├── user_repo.go       # User data access
│   │   ├── api_key_repo.go    # API key data access
│   │   └── usage_repo.go      # Usage data access
│   └── service/
│       ├── auth_service.go    # Auth business logic
│       ├── proxy_service.go   # Proxy business logic
│       ├── usage_service.go   # Usage calculations
│       └── billing_service.go # Billing calculations
├── pkg/
│   ├── aiprovider/
│   │   ├── provider.go        # Provider interface
│   │   ├── openai.go        # OpenAI adapter
│   │   └── anthropic.go      # Anthropic adapter
│   └── utils/
│       ├── hash.go           # Hashing utilities
│       └── response.go       # HTTP response helpers
├── api/
│   └── openapi.yaml         # OpenAPI specification (future)
├── docs/                     # MkDocs documentation
├── .github/
│   └── workflows/
│       ├── ci.yml           # CI pipeline
│       └── cd.yml           # Deploy pipeline
├── .golangci.yml            # Go linter config
├── Makefile                 # Build commands
├── config.yaml              # Default config
├── mkdocs.yml              # MkDocs config
├── go.mod
└── go.sum
```

## Directory Purpose

| Directory | Purpose |
|-----------|---------|
| `cmd/` | Application entry points |
| `internal/` | Private application code |
| `pkg/` | Public packages importable by external projects |
| `api/` | API specifications (OpenAPI) |
| `docs/` | MkDocs documentation source |
| `.github/workflows/` | GitHub Actions CI/CD |

## Naming Conventions

- **Files** — lowercase with underscores: `api_key_repo.go`
- **Types** — PascalCase: `APIKey`, `UsageSummary`
- **Interfaces** — PascalCase with `er` suffix: `Provider`, `Repository`
- **Tests** — `_test.go` suffix: `auth_service_test.go`
- **Migrations** — numbered with up/down: `000001_init_schema.up.sql`
