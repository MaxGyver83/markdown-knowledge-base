# Docker

Docker is used to containerize all three components of the Markdown Knowledge Base.

## Multi-Service Architecture

The project runs three Docker containers orchestrated by docker-compose:

| Service | Image | Port | Description |
|---------|-------|------|-------------|
| postgres | postgres:16-alpine | 5432 | Database |
| backend | markdown-api (custom) | 8080 | Go REST API |
| frontend | nginx:alpine | 8085 | Static file server + reverse proxy |

## docker-compose.yml

The `infra/docker-compose.yml` file defines the three services:

```yaml
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: markdownkb
      POSTGRES_PASSWORD: postgres
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 3s
      retries: 5

  backend:
    build: ../backend/markdown-api
    environment:
      DATABASE_URL: postgres://postgres:postgres@postgres:5432/markdownkb?sslmode=disable
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - markdown-data:/app/data/documents

  frontend:
    image: nginx:alpine
    ports:
      - "8085:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf:ro
    depends_on:
      - backend

volumes:
  pgdata:
  markdown-data:
```

## Multi-Stage Backend Build

The backend Dockerfile uses a two-stage build to keep the final image small:

```dockerfile
# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /server
EXPOSE 8080
CMD ["/server"]
```

The final image is approximately 15 MB and contains only the compiled binary.

## Running Locally

```sh
cd infra
docker compose up -d       # Start all services
docker compose down         # Stop all services
docker compose build backend # Rebuild the backend image
docker compose logs -f      # Follow logs
```
