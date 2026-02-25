# Repository Guidelines

## Project Structure & Module Organization
- `main.go` is the entry point and wires the Cobra CLI.
- `app/` contains application logic (commands, HTTP controllers, middleware, models, requests, policies).
- `routes/` defines API routing.
- `bootstrap/` initializes logger, database, cache, Redis, and routes.
- `config/` holds configuration loaders and defaults.
- `database/` contains migrations, seeders, and factories.
- `pkg/` contains reusable packages (auth, cache, jwt, logger, paginator, etc.).
- `storage/` is for runtime artifacts like logs.

## Build, Test, and Development Commands
- `go run main.go` runs the default `serve` command (starts the API server).
- `go run main.go serve` explicitly starts the web server.
- `go run main.go migrate up|down|reset|refresh|fresh` manages database migrations.
- `go run main.go seed [SeederName]` seeds data (all seeders if no name).
- `go run main.go key` generates and prints an app key.
- `go run main.go -h` lists all supported CLI commands.
- `go build ./...` compiles all packages.

## Coding Style & Naming Conventions
- Use standard Go formatting (`gofmt`); tabs for indentation.
- Package names are lowercase and short (for example `pkg/jwt`, `app/requests`).
- Files follow snake case in existing areas (for example `topic_request.go`).
- Keep API handlers in `app/http/controllers/api/v1` and requests in `app/requests`.

## Testing Guidelines
- No tests are currently present. New behavior should include `_test.go` files.
- Run `go test ./...` locally to execute all tests.

## Commit & Pull Request Guidelines
- Commit messages are short, imperative, and scoped (for example "Remove redundant parameters").
- PRs should include:
  - A concise summary of changes.
  - Related issue links if applicable.
  - Migration/seed instructions when data shape changes.
  - Test results or a note if tests are not available.

## Configuration Tips
- Use `.env.example` as the baseline for local configuration.
- Set `--env=testing` to load `.env.testing` (or similar) via the CLI flag.
- Ensure database and Redis settings match your local services.

## Dependency Status Checks
- `github.com/redis/go-redis/v9` is the maintained Redis client; the older `github.com/go-redis/redis` repo redirects to it.
- `github.com/golang-jwt/jwt` is the maintained JWT library; `dgrijalva/jwt-go` is archived.
- `github.com/go-faker/faker` notes it was previously `bxcodec/faker`.
- No archived banner found for: `gin-gonic/gin`, `gorm.io/gorm`, `spf13/viper`, `spf13/cobra`, `ulule/limiter`, `mojocn/base64Captcha`, `thedevsaddam/govalidator`, `mgutz/ansi`, `jordan-wright/email`, `gertd/go-pluralize`, `iancoleman/strcase`.
