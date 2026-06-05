# AI Budget Splitter Backend

A Go-based backend service that acts as an **API key rerouter**, enabling organizations (up to 10 members) to track and split AI usage costs.

## Overview

This service solves the problem of dividing AI API bills among team members. Instead of each member having their own API key, the backend:

1. Issues **generated API keys** to each organization member
2. **Proxies requests** to AI providers (OpenAI, Anthropic) using shared real API keys
3. **Tracks usage** per user (tokens, cost)
4. **Calculates billing splits** based on actual usage percentages

## Features

- **API Key Rerouting** — Users call the backend with their generated key; the backend proxies to AI providers using real keys
- **Organization Management** — Manage organizations with up to 10 members per org
- **Usage Tracking** — Per-user tracking of input/output tokens and cost
- **Billing Split** — Automatic percentage-based cost allocation
- **Multi-Provider Support** — OpenAI and Anthropic API proxies
- **Rate Limiting** — Per-user rate limiting middleware
- **JWT Authentication** — Secure token-based auth for backend operations

## Quick Links

- [Installation](getting-started/installation.md)
- [Configuration](getting-started/configuration.md)
- [API Reference](api/auth.md)
- [Architecture](architecture/overview.md)

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.22+ |
| Web Framework | Gin |
| Database | PostgreSQL + pgx |
| Configuration | Viper |
| Authentication | JWT + bcrypt |
| CI/CD | GitHub Actions |

## Project Status

This is an initial implementation. See the [GitHub repository](https://github.com/jackdafox/backend-ai-budget-splitter) for current progress.
