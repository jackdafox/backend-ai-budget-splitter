# Organization API

All endpoints require JWT authentication via `Authorization: Bearer <token>` header.

## Get Organization

**Endpoint:** `GET /api/v1/org`

```bash
curl http://localhost:8080/api/v1/org \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "My Team",
  "created_at": "2024-06-01T00:00:00Z",
  "updated_at": "2024-06-01T00:00:00Z"
}
```

## Update Organization

**Endpoint:** `PUT /api/v1/org`

```bash
curl -X PUT http://localhost:8080/api/v1/org \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_HERE" \
  -d '{
    "name": "New Team Name"
  }'
```

**Response** `200 OK`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "New Team Name",
  "created_at": "2024-06-01T00:00:00Z",
  "updated_at": "2024-06-01T12:00:00Z"
}
```

## List Members

**Endpoint:** `GET /api/v1/org/members`

```bash
curl http://localhost:8080/api/v1/org/members \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "members": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "org_id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2024-06-01T00:00:00Z"
    }
  ]
}
```

## Add Member

Maximum 10 members per organization.

**Endpoint:** `POST /api/v1/org/members`

```bash
curl -X POST http://localhost:8080/api/v1/org/members \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_HERE" \
  -d '{
    "email": "member@example.com",
    "password": "userpassword123",
    "role": "member"
  }'
```

**Response** `201 Created`

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "org_id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "member@example.com",
  "role": "member",
  "created_at": "2024-06-01T12:00:00Z"
}
```

## Remove Member

**Endpoint:** `DELETE /api/v1/org/members/:user_id`

```bash
curl -X DELETE http://localhost:8080/api/v1/org/members/550e8400-e29b-41d4-a716-446655440003 \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "message": "member removed"
}
```

## Error Responses

| Status | Error | Cause |
|--------|-------|-------|
| `400` | `organization member limit reached (max 10)` | Member cap exceeded |
| `400` | `email already registered` | Email in use |
| `404` | `user not found` | User ID doesn't exist |
