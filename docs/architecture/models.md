# Data Models

## Organization

Represents a team or company that shares AI API costs.

```go
type Organization struct {
    ID        uuid.UUID  `json:"id"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}
```

## User

Belongs to an organization and has a role.

```go
type MemberRole string

const (
    RoleAdmin  MemberRole = "admin"
    RoleMember MemberRole = "member"
)

type User struct {
    ID             uuid.UUID `json:"id"`
    OrganizationID uuid.UUID `json:"org_id"`
    Email          string    `json:"email"`
    PasswordHash   string    `json:"-"`        // Never exposed
    Role           MemberRole `json:"role"`
    CreatedAt      time.Time `json:"created_at"`
}
```

## APIKey

Generated for users to authenticate proxy requests. Stored as bcrypt hash.

```go
type APIKey struct {
    ID             uuid.UUID  `json:"id"`
    UserID         uuid.UUID  `json:"user_id"`
    OrganizationID uuid.UUID  `json:"org_id"`
    KeyHash        string     `json:"-"`           // bcrypt hash
    KeyPrefix      string     `json:"key_prefix"`  // First 8 chars
    Name           string     `json:"name"`
    LastUsedAt     *time.Time `json:"last_used_at,omitempty"`
    ExpiresAt      *time.Time `json:"expires_at,omitempty"`
    CreatedAt      time.Time  `json:"created_at"`
}
```

## UsageRecord

Tracks each AI API call for billing attribution.

```go
type UsageRecord struct {
    ID             uuid.UUID `json:"id"`
    OrganizationID uuid.UUID `json:"org_id"`
    UserID         uuid.UUID `json:"user_id"`
    Provider       string    `json:"provider"` // "openai" or "anthropic"
    Model          string    `json:"model"`
    InputTokens    int       `json:"input_tokens"`
    OutputTokens   int       `json:"output_tokens"`
    Cost           float64   `json:"cost"`
    RecordedAt     time.Time `json:"recorded_at"`
}
```

## UsageSummary

Aggregated usage per user.

```go
type UsageSummary struct {
    UserID         uuid.UUID `json:"user_id"`
    TotalInput     int       `json:"total_input_tokens"`
    TotalOutput    int       `json:"total_output_tokens"`
    TotalCost      float64   `json:"total_cost"`
    UsagePercent   float64   `json:"usage_percent"` // % of org total
}
```

## OrgUsageSummary

Aggregated usage for entire organization.

```go
type OrgUsageSummary struct {
    OrganizationID uuid.UUID      `json:"org_id"`
    TotalCost      float64        `json:"total_cost"`
    TotalInput     int            `json:"total_input_tokens"`
    TotalOutput    int            `json:"total_output_tokens"`
    UserSummaries  []UsageSummary `json:"user_summaries"`
}
```

## BillingBreakdown

Detailed billing per member.

```go
type BillingBreakdown struct {
    OrgID         uuid.UUID        `json:"org_id"`
    PeriodStart   time.Time        `json:"period_start"`
    PeriodEnd     time.Time        `json:"period_end"`
    TotalCost     float64          `json:"total_cost"`
    TotalInput    int              `json:"total_input_tokens"`
    TotalOutput   int              `json:"total_output_tokens"`
    MemberBillets []MemberBillable `json:"member_billets"`
}

type MemberBillable struct {
    UserID       uuid.UUID `json:"user_id"`
    Email        string    `json:"email"`
    Role         string    `json:"role"`
    InputTokens  int       `json:"input_tokens"`
    OutputTokens int       `json:"output_tokens"`
    Cost         float64   `json:"cost"`
    Percentage   float64   `json:"percentage"`
    OwedAmount   float64   `json:"owed_amount"`
}
```
