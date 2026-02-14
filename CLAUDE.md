# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run the application
go run Main.go

# Build executable
go build -o DigiPassAuthenticationApi.exe

# Format code
go fmt ./...
```

No test suite is set up yet. No linter configuration exists beyond `go fmt`.

## Required Environment

- Go 1.25.6
- PostgreSQL database
- `.env` file for local development with: `DATABASE_URL`, `JWT_SECRET`, `JWT_ALGO`
- On Railway, env vars are set in the platform and `.env` is not loaded

## Architecture

This is a Go REST API using **Echo v5** (web framework) and **GORM** (ORM) with PostgreSQL.

### Layered Structure

**Entry point:** `Main.go` loads env vars, `App.go` initializes Echo, DB, middleware, and routes.

**Layers (request flow):**

1. **Routes** (`routes/v1/`) — Register endpoint groups under `/v1`. Each entity has its own route file that calls `RegisterXxxRoutes(group)`.
2. **Handlers** (`handlers/`) — Bind/validate request bodies, extract DB from context, wrap service calls in transactions, map service errors to HTTP status codes.
3. **Services** (`services/`) — Business logic. Constructed per-request with a DB/tx reference: `services.NewXxxService(tx)`. Custom errors defined in `services/Errors.go` and checked with `errors.Is()`.
4. **Models** (`packages/models/Models.go`) — All GORM models in a single file. UUID primary keys, JSON tags, validation tags (`go-playground/validator`), and GORM relationship definitions.

### Key Patterns

- **DB injection:** GORM `*gorm.DB` is set on Echo context via middleware (`c.Set("db", db)`), retrieved in handlers with `c.Get("db").(*gorm.DB)`.
- **Transactions:** `utils.WithTransaction(db, func(tx *gorm.DB) error { ... })` wraps service calls with automatic rollback on error.
- **Validation:** Echo's `Validator` is set to a `ValidatorWrapper` around `go-playground/validator`. Handlers call `c.Bind()` then `c.Validate()` on inline request structs.
- **JWT:** Custom HMAC-SHA256 implementation in `packages/jwt/Create.go` (not a third-party library).
- **Error responses:** JSON format `{"error": "message"}`. Handlers map known service errors (e.g., `ErrAccountAlreadyExists`) to appropriate HTTP status codes.

### Entity Model

Core entities: Account → Tenant (auto-created with account), Client, User, AccountUser (admin/staff roles: owner/admin/member), Session, MagicLink, OAuth tokens (AuthorizationCode, AccessToken, RefreshToken, IDToken), AuditLog, UserConsent.

## Conventions

- PascalCase file names (e.g., `Account.go`, `AccountUsers.go`)
- Constructor functions: `NewXxxHandler()`, `NewXxxService(tx)`
- Handler methods are receiver methods on handler structs
- Route registration functions: `RegisterXxxRoutes(group)`
- Server runs on port 1323
