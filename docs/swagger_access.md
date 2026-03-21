# Swagger Access Guide

## Prerequisites

- `.env` is configured
- PostgreSQL is running and reachable

## Start the API

```bash
go run cmd/run/main.go
```

By default, server runs on `http://localhost:8080`.

## Open Swagger UI

Use this URL in your browser:

- `http://localhost:8080/swagger/index.html`

Raw OpenAPI JSON:

- `http://localhost:8080/swagger/doc.json`

## Regenerate docs after API changes

If you add/update Swagger annotations in code:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
"$(go env GOPATH)/bin/swag" init -g cmd/run/main.go -o docs/swagger
```

## How to use secured endpoints

1. Login using `POST /api/v1/auth/login`.
2. Copy the returned JWT token.
3. Click **Authorize** in Swagger UI.
4. Enter:
   - `Bearer <your-token>`

## Regenerate docs after API changes

If you add/update Swagger annotations in code:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
"$(go env GOPATH)/bin/swag" init -g cmd/run/main.go -o docs/swagger
```

