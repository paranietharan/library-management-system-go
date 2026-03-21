# Progress Tracker

## Connection Details

TBD (no DB/VM connection details provided in docs)

## Session Log

### 2026-03-20
#### Repo Inspection
- [x] Read `docs/architecture.md` to extract intended endpoints/roles.
- [x] Inspected existing Go implementation: `auth` + `books` + `reviews` + `comments` routes/services/repos.
- [x] Verified DB migration file `internal/database/migration/1_init.up.sql` defines additional tables (`lendings`, `reservations`, `fines`, `complaints`, `articles`).
- [x] Ran `go test ./...` (completed successfully; `exit_code: 0`).

#### Implementation
- [x] Fixed auth context key mismatch by changing `c.GetUint("userID")` to `c.GetUint("user_id")` in `internal/handler/review_handler.go` and `internal/handler/comment_handler.go`, and updated unit tests accordingly.
- [x] Fixed Gin wildcard route conflict panic by normalizing nested route params to use `:id` under `/books/:id/*` and `/articles/:id/*`; updated related handlers/tests.
- [x] Added Swagger support with route `GET /swagger/*any` in `cmd/api/router.go`.
- [x] Added Swagger metadata annotations in `cmd/run/main.go` and generated OpenAPI files under `docs/swagger/`.
- [x] Added documentation for Swagger usage in `docs/swagger_access.md` and updated `README.md` with access URLs and regeneration steps.
- [x] Resolved Swagger generator/library version mismatch by upgrading `github.com/swaggo/swag` to `v1.16.6`.
- [x] Added SQL migration `internal/database/migration/2_add_book_reviews_comments_and_article_subresources.up.sql` to create missing tables: `reviews`, `comments`, `article_reviews`, `article_comments`, `article_ratings` (plus indexes/triggers).
- [x] Implemented `articles` endpoints and article sub-resources in code (`internal/domain`, `internal/dto`, `internal/repository`, `internal/service`, `internal/handler`) and wired routes in `cmd/api/routes.go`.
- [x] Implemented `lendings`, `reservations`, `fines`, and `complaints` endpoints and wired routes in `cmd/api/routes.go`.
- [x] Implemented admin-only `users` and `roles` endpoints (static role enum mapping) and wired routes in `cmd/api/routes.go`.
- [x] Ran `go test ./...` again after router/service wiring (completed successfully; `exit_code: 0`).

#### Scope & Next Steps (Pending)
- [x] Fix API context key mismatch in handlers (`middleware` sets `user_id`, but `review_handler`/`comment_handler` read `userID`), then update related unit tests.
- [x] Implement missing domain/models/services/handlers/routes for `articles` (and article sub-resources), `lendings`, `reservations`, `fines`, and `complaints`.
- [x] Implement admin endpoints for `users` and `roles` from `docs/architecture.md`.
- [x] Extend DB migration to add missing tables required by existing code (`reviews`, `comments`) and by article sub-resources (`article_reviews`, `article_comments`, `article_ratings`).