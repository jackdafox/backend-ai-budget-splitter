# Deployment

## Docker

### Dockerfile

```dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/server .
COPY config.yaml .

EXPOSE 8080
CMD ["./server"]
```

### Build and Run

```bash
docker build -t ai-budget-splitter .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db?sslmode=disable" \
  -e OPENAI_API_KEY="sk-..." \
  -e ANTHROPIC_API_KEY="sk-ant-..." \
  -e JWT_SECRET="your-secret" \
  ai-budget-splitter
```

## Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://postgres:postgres@db:5432/ai_budget_splitter?sslmode=disable
      - OPENAI_API_KEY=${OPENAI_API_KEY}
      - ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      db:
        condition: service_healthy

  db:
    image: postgres:16
    environment:
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: ai_budget_splitter
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

```bash
docker compose up -d
```

## Kubernetes

### Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ai-budget-splitter
spec:
  replicas: 2
  selector:
    matchLabels:
      app: ai-budget-splitter
  template:
    metadata:
      labels:
        app: ai-budget-splitter
    spec:
      containers:
        - name: server
          image: your-registry/ai-budget-splitter:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: ai-budget-splitter-secrets
            - name: OPENAI_API_KEY
              valueFrom:
                secretKeyRef:
                  name: ai-budget-splitter-secrets
            - name: ANTHROPIC_API_KEY
              valueFrom:
                secretKeyRef:
                  name: ai-budget-splitter-secrets
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: ai-budget-splitter-secrets
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "512Mi"
              cpu: "500m"
```

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: ai-budget-splitter
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: ai-budget-splitter
```

## Production Checklist

- [ ] PostgreSQL with SSL enabled
- [ ] Strong JWT secret (32+ random characters)
- [ ] API keys stored in secrets manager
- [ ] HTTPS via reverse proxy (nginx/Caddy)
- [ ] Rate limiting configured appropriately
- [ ] Logging to stdout (for log aggregation)
- [ ] Health check endpoint (`/health`)
- [ ] Graceful shutdown configured
- [ ] Database connection pooling tuned
- [ ] Monitoring (Prometheus metrics — future)
