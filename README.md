# Gohub

A Gin-based forum API written in Go. This repository provides a CLI
entrypoint for running the API server, migrations, seeders, and scaffolding.

## Requirements

- Go 1.26
- Redis for non-testing runtime cache, captcha, verify-code, and limiter storage
- A supported database (PostgreSQL/MySQL/SQLite)
- Python 3.13+ for the real HTTP API smoke script

## Quick Start

```bash
cp .env.example .env
# edit .env with your database/redis credentials

go run main.go key
go run main.go migrate up
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

# real HTTP API smoke test against a running server
python3 scripts/api_smoke_test.py --base-url http://127.0.0.1:3000
```

Command initialization is scoped by command:

- `key` and `make` do not initialize DB or Redis.
- `migrate` and `seed` initialize the logger and DB only.
- `cache` initializes the logger, Redis, and cache only.
- `serve` initializes the full API runtime.

## Configuration

- Use `.env.example` as the baseline.
- `--env=testing` loads `.env.testing` (if present).
- `APP_ENV_PATH` points to an explicit env file path (useful for tests).
  It takes precedence over `--env` and the default `.env`.
- In tests, setting `CONSOLE_SILENT=1` silences console output when
  `APP_ENV=testing`.
- `.env.testing` uses SQLite at `/tmp/gohub_api_smoke.db` and is intended for
  local smoke tests.

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

Common response codes:

- `OK`
- `CREATED`
- `ERR_BAD_REQUEST`
- `ERR_UNAUTHORIZED`
- `ERR_FORBIDDEN`
- `ERR_NOT_FOUND`
- `ERR_VALIDATION`
- `ERR_UNPROCESSABLE`
- `ERR_INTERNAL`

Pagination query parameters:

- `offset`: integer, minimum `0`
- `limit`: integer, `1..100`
- `sort`: one of `id`, `created_at`, `updated_at`
- `order`: one of `asc`, `desc`

## API Routes

All routes below use the `/api/v1` prefix unless `API_DOMAIN` is configured,
in which case the prefix is `/v1`.

| Method | Path | Auth | Notes |
|---|---|---:|---|
| POST | `/auth/signup/phone/exist` | No | Check whether a phone is registered |
| POST | `/auth/signup/email/exist` | No | Check whether an email is registered |
| POST | `/auth/signup/using-phone` | No | Sign up with phone verification code |
| POST | `/auth/signup/using-email` | No | Sign up with email verification code |
| POST | `/auth/verify-codes/captcha` | No | Generate image captcha |
| POST | `/auth/verify-codes/phone` | No | Send phone verification code |
| POST | `/auth/verify-codes/email` | No | Send email verification code |
| POST | `/auth/login/using-phone` | No | Login with phone verification code |
| POST | `/auth/login/using-password` | No | Login with username/email/phone and password |
| POST | `/auth/login/refresh-token` | Bearer | Refresh access token |
| POST | `/auth/password-reset/using-phone` | No | Reset password by phone verification code |
| POST | `/auth/password-reset/using-email` | No | Reset password by email verification code |
| GET | `/user` | Bearer | Current user |
| GET | `/users` | No | Paginated users |
| PUT | `/users` | Bearer | Update profile |
| PUT | `/users/email` | Bearer | Update email |
| PUT | `/users/phone` | Bearer | Update phone |
| PUT | `/users/password` | Bearer | Update password |
| PUT | `/users/avatar` | Bearer | Upload avatar as multipart field `avatar` |
| GET | `/categories` | No | Paginated categories |
| POST | `/categories` | Bearer | Create category |
| PUT | `/categories/:id` | Bearer | Update owned category |
| DELETE | `/categories/:id` | Bearer | Delete owned category |
| GET | `/topics` | No | Paginated topics |
| GET | `/topics/:id` | No | Show topic |
| POST | `/topics` | Bearer | Create topic |
| PUT | `/topics/:id` | Bearer | Update owned topic |
| DELETE | `/topics/:id` | Bearer | Delete owned topic |
| GET | `/links` | No | List links |

## Real API Smoke Test

The smoke script exercises every route registered in `routes/api.go` through
real HTTP requests.

```bash
go run main.go --env=testing migrate fresh
go run main.go --env=testing serve
python3 scripts/api_smoke_test.py --base-url http://127.0.0.1:3000
```

In `APP_ENV=testing`, captcha and verification-code checks support the
built-in test shortcuts:

- captcha ID: `captcha_skip_test`
- phone prefix: `000`
- email suffix: `@testing.com`

## Notes

- Use `go test ./...` to run tests and validate behavior.
- Use `go vet ./...` before committing Go changes.
- Chinese documentation is in `README-zh.md`.
