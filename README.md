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

## Notes
- There are currently no test files; use `go test ./...` to validate compilation.
- Chinese documentation is in `README-zh.md`.
