# Usage & Billing API

## Get Current Usage

Get usage for the current billing period (month to date).

**Endpoint:** `GET /api/v1/usage`

```bash
curl http://localhost:8080/api/v1/usage \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "org_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_cost": 12.50,
  "total_input_tokens": 50000,
  "total_output_tokens": 30000,
  "user_summaries": [
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "total_input_tokens": 30000,
      "total_output_tokens": 20000,
      "total_cost": 7.50,
      "usage_percent": 60.0
    },
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440003",
      "total_input_tokens": 20000,
      "total_output_tokens": 10000,
      "total_cost": 5.00,
      "usage_percent": 40.0
    }
  ]
}
```

## Get Usage History

Get usage for a specific time period.

**Endpoint:** `GET /api/v1/usage/history`

```bash
curl "http://localhost:8080/api/v1/usage/history?start=2024-01-01T00:00:00Z&end=2024-01-31T23:59:59Z" \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Query Parameters:**

| Parameter | Format | Description |
|-----------|--------|-------------|
| `start` | RFC3339 | Start of period (default: first of month) |
| `end` | RFC3339 | End of period (default: now) |

## Get Billing Breakdown

Get detailed billing per member for a period.

**Endpoint:** `GET /api/v1/billing`

```bash
curl "http://localhost:8080/api/v1/billing?start=2024-01-01T00:00:00Z&end=2024-01-31T23:59:59Z" \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

```json
{
  "org_id": "550e8400-e29b-41d4-a716-446655440000",
  "period_start": "2024-01-01T00:00:00Z",
  "period_end": "2024-01-31T23:59:59Z",
  "total_cost": 100.00,
  "total_input_tokens": 400000,
  "total_output_tokens": 200000,
  "member_billets": [
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440001",
      "email": "admin@example.com",
      "role": "admin",
      "input_tokens": 240000,
      "output_tokens": 120000,
      "cost": 60.00,
      "percentage": 60.0,
      "owed_amount": 60.00
    },
    {
      "user_id": "550e8400-e29b-41d4-a716-446655440003",
      "email": "member@example.com",
      "role": "member",
      "input_tokens": 160000,
      "output_tokens": 80000,
      "cost": 40.00,
      "percentage": 40.0,
      "owed_amount": 40.00
    }
  ]
}
```

## Export Billing Data

Export billing data for a period.

**Endpoint:** `GET /api/v1/billing/export`

```bash
curl "http://localhost:8080/api/v1/billing/export?start=2024-01-01T00:00:00Z&end=2024-01-31T23:59:59Z" \
  -H "Authorization: Bearer YOUR_JWT_HERE"
```

**Response** `200 OK`

Returns JSON data (CSV/PDF export formats planned).

## Cost Calculation

Costs are approximate and based on published AI provider pricing:

| Provider | Input | Output |
|----------|-------|--------|
| OpenAI (GPT-3.5) | $0.001/1K tokens | $0.003/1K tokens |
| OpenAI (GPT-4) | $0.01-0.03/1K tokens | $0.03-0.06/1K tokens |
| Anthropic (Claude 3) | $0.015/1K tokens | $0.075/1K tokens |

> Actual costs may vary based on model and provider updates.

## Error Responses

| Status | Error | Cause |
|--------|-------|-------|
| `400` | `invalid start date format` | Malformed RFC3339 date |
| `400` | `invalid end date format` | Malformed RFC3339 date |
