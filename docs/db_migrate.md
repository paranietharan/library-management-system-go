# Steps to Migrate Database

## Run migrations
```bash
go run ./cmd/db_migrate migrate
```

## Rollback last migration
```bash
go run ./cmd/db_migrate rollback
```

## Seed default users
```bash
go run ./cmd/db_migrate seed
```