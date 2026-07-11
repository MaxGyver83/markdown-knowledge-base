# Markdown Knowledge Base

- Backend: REST API written in Go (using go-chi for routing)
- Frontend: HTML + Javascript + CSS
- Database: PostgreSQL
- CI/CD: GitHub Actions

## Build and run with Docker

```sh
cd infra
docker compose up -d
```

### Stop

```sh
cd infra
docker compose down
```

## Rebuild

Backend for example:

```sh
cd infra
docker compose build backend
```
