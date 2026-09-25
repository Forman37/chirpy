# Chirpy

Chirpy is a small REST API for creating users and short posts ("chirps"). It uses PostgreSQL for persistence, SQLC for type-safe database queries, Goose for migrations, Argon2id password hashing, JWT access tokens, refresh tokens, and a Polka webhook for Chirpy Red upgrades.

## Requirements

- Go 1.25.1 or newer
- PostgreSQL
- [Goose](https://github.com/pressly/goose) for database migrations
- [SQLC](https://sqlc.dev/) when changing SQL queries or schema definitions

Install the command-line tools with Go:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Ensure Go's binary directory (usually `$(go env GOPATH)/bin`) is on your `PATH`.

## Setup

1. Create a PostgreSQL database:

   ```sh
   createdb chirpy
   ```

2. Create a `.env` file in the repository root. The application loads it automatically when it starts:

   ```dotenv
   DB_URL=postgres://postgres:password@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWTSECRET=replace-with-a-long-random-secret
   POLKA_KEY=replace-with-the-polka-webhook-api-key
   ```

   `PLATFORM=dev` enables `POST /admin/reset`; use a non-`dev` value outside local development. Keep `.env` private: it is intentionally ignored by Git.

3. Apply migrations from the repository root:

   ```sh
   goose -dir sql/schema postgres "$DB_URL" up
   ```

   If `DB_URL` only exists in `.env`, either set it in your shell or provide the connection string directly:

   ```sh
   goose -dir sql/schema postgres "postgres://postgres:password@localhost:5432/chirpy?sslmode=disable" up
   ```

4. Start the API:

   ```sh
   go run .
   ```

   The server listens at `http://localhost:8080`.

## Common Commands

Run all commands below from the repository root. The migration and SQLC commands originate from `commands.txt`; `-dir sql/schema` is required because migrations are stored there rather than in Goose's default directory.

| Task | Command |
| --- | --- |
| Run the API | `go run .` |
| Run tests | `go test ./...` |
| Format Go code | `go fmt ./...` |
| Apply all pending migrations | `goose -dir sql/schema postgres "$DB_URL" up` |
| Revert the most recent migration | `goose -dir sql/schema postgres "$DB_URL" down` |
| Show migration state | `goose -dir sql/schema postgres "$DB_URL" status` |
| Create a SQL migration | `goose -dir sql/schema create describe_change sql` |
