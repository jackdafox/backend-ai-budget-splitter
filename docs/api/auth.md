# Authentication API

## Register Organization

Create a new organization with an admin user.

**Endpoint:** `POST /api/v1/auth/register`

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "org_name": "My Team",
    "email": "admin@example.com",
    "password": "securepassword123"
  }'
```

**Response** `201 Created`

```json
{
  "org_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440001",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Login

Authenticate and receive a JWT.

**Endpoint:** `POST /api/v1/auth/login`

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "securepassword123"
  }'
```

**Response** `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "org_id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "admin@example.com",
    "role": "admin"
  }
}
```

## Generate API Key

Create a new API key for a user. The full key is returned only once.

**Endpoint:** `POST /api/v1/auth/keys`

```bash
curl -X POST http://localhost:8080/api/v1/auth/keys \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_HERE" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "My API Key"
  }'
```

**Response** `201 Created`

```json
{
  "api_key": "f47ac10b-58cc-4372-a567-0e02d2dbf1e2...",
  "key_prefix": "f47ac10b",
  "key_id": "550e8400-e29b-41d4-a716-446655440002",
  "name": "My API Key",
  "expires_at": "2024-09-01T00:00:00Z",
  "created_at": "2024-06-01T00:00:00Z"
}
```

> **Important:** Save the `api_key` — it is shown only once.

## List API Keys

Get all API keys for the authenticated user.

**Endpoint:** `GET /api/v1/auth/keys`

```bash
curl http://localhost:8080/api/v1/auth/keys \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "api_keys": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "org_id": "550e8400-e29b-41d4-a716-446655440000",
      "key_prefix": "f47ac10b",
      "name": "My API Key",
      "last_used_at": "2024-06-01T12:00:00Z",
      "expires_at": "2024-09-01T00:00:00Z",
      "created_at": "2024-06-01T00:00:00Z"
    }
  ]
}
```

## Revoke API Key

Delete an API key.

**Endpoint:** `DELETE /api/v1/auth/keys/:key_id`

```bash
curl -X DELETE http://localhost:8080/api/v1/auth/keys/550e8400-e29b-41d4-a716-446655440002 \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "message": "API key revoked"
}
```

## Error Responses

| Status | Error | Cause |
|--------|-------|-------|
| `400` | `email already registered` | Email in use |
| `401` | `invalid credentials` | Wrong password |
| `401` | `missing authorization header` | No JWT provided |
| `404` | `user not found` | User ID doesn't exist |
