# Installation

## Prerequisites

- **Go 1.22+** — [Install Go](https://go.dev/doc/install)
- **PostgreSQL 16+** — [Install PostgreSQL](https://www.postgresql.org/download/)
- **Make** — Usually pre-installed on Linux/macOS; on Windows use WSL or install via scoop/chocolatey

## Clone the Repository

```bash
git clone https://github.com/jackdafox/backend-ai-budget-splitter.git
cd backend-ai-budget-splitter
```

## Install Dependencies

```bash
make deps
```

This runs `go mod download` and `go mod tidy` to fetch all dependencies.

## Database Setup

1. Create a PostgreSQL database:

```sql
CREATE DATABASE ai_budget_splitter;
```

2. Note your connection string:

```
postgres://postgres:password@localhost:5432/ai_budget_splitter?sslmode=disable
```

3. Run the migrations (see [Configuration](configuration.md) for setup):

```bash
# After configuring config.yaml
make migrate
```

## Configuration

Copy the example config and customize:

```bash
cp config.yaml config.local.yaml
```

Edit `config.local.yaml` with your settings. See [Configuration](configuration.md) for details.

## Build

```bash
make build
```

This produces a binary at `bin/server`.

## Run

```bash
make run
```

Or run the binary directly:

```bash
./bin/server
```

The server starts on `http://localhost:8080`.

## Verify

```bash
curl http://localhost:8080/health
```

Expected response:

```json
{"status":"ok"}
```
