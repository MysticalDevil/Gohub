# Gohub

A Gin-based forum API written in Go. This repository provides a CLI entrypoint for running the API server, migrations, seeders, and scaffolding.

## Requirements
- Go 1.26
- Redis (for cache, captcha, limiter)
- A supported database (PostgreSQL/MySQL/SQLite)

## Quick Start
```bash
cp .env.example .env
# edit .env with your database/redis credentials

go run main.go serve
```

## Project Layout
- `main.go` entrypoint (Cobra CLI)
- `app/` application logic (commands, controllers, middleware, models, requests)
- `routes/` API routing
- `bootstrap/` initialization (logger, DB, Redis, cache, routes)
- `config/` configuration loading and defaults
- `database/` migrations, seeders, factories
- `pkg/` shared packages (auth, cache, jwt, logger, paginator, etc.)
- `storage/` runtime artifacts (logs)

## Common Commands
```bash
# list CLI commands
go run main.go -h

# start the API server
go run main.go serve

# database migrations
go run main.go migrate up|down|reset|refresh|fresh

# seed data (all or by name)
go run main.go seed
# go run main.go seed UsersSeeder

# generate app key
go run main.go key
```

## Configuration
- Use `.env.example` as the baseline.
- `--env=testing` loads `.env.testing` (if present).
- `APP_ENV_PATH` points to an explicit env file path (useful for tests). It takes precedence over `--env` and the default `.env`.
- In tests, setting `CONSOLE_SILENT=1` silences console output when `APP_ENV=testing`.

## API Responses
All API responses use a standard envelope:
```json
{
  "code": "OK",
  "msg": "OK",
  "data": {}
}
```

Validation or processing errors return:
```json
{
  "code": "ERR_VALIDATION",
  "msg": "Request verification failed, please see errors for details",
  "errors": {
    "field": ["message"]
  }
}
```

Pagination uses `offset/limit` and returns:
```json
{
  "code": "OK",
  "msg": "OK",
  "data": {
    "items": [],
    "offset": 0,
    "limit": 20,
    "total": 200
  }
}
```

## Notes
- Use `go test ./...` to run tests and validate behavior.
- Chinese documentation is in `README-zh.md`.
