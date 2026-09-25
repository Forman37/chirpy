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

Run all commands below from the repository root.

| Task | Command |
| --- | --- |
| Run the API | `go run .` |
| Run tests | `go test ./...` |
| Format Go code | `go fmt ./...` |
| Apply all pending migrations | `goose -dir sql/schema postgres "$DB_URL" up` |
| Revert the most recent migration | `goose -dir sql/schema postgres "$DB_URL" down` |
| Show migration state | `goose -dir sql/schema postgres "$DB_URL" status` |
| Create a SQL migration | `goose -dir sql/schema create describe_change sql` |
| Regenerate SQLC code | `sqlc generate` |

After editing a file in `sql/schema` or `sql/queries`, run `sqlc generate` to update `internal/database` before testing or committing.

## API

All JSON endpoints accept and return `application/json` unless noted otherwise. Protected endpoints require an `Authorization: Bearer <token>` header.

Chirps are limited to 140 characters. `kerfuffle`, `sharbert`, and `fornax` are replaced with `****`.

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `GET` | `/api/healthz` | No | Returns plain-text `OK`. |
| `POST` | `/api/validate_chirp` | No | Validates and cleans a chirp body without storing it. |
| `POST` | `/api/users` | No | Creates a user. |
| `PUT` | `/api/users` | JWT | Updates the authenticated user's email and password. |
| `POST` | `/api/login` | No | Authenticates a user and returns access and refresh tokens. |
| `POST` | `/api/refresh` | Refresh token | Returns a new access token. |
| `POST` | `/api/revoke` | Refresh token | Revokes a refresh token. |
| `GET` | `/api/chirps` | No | Lists all chirps. |
| `GET` | `/api/chirps/{chirpID}` | No | Retrieves one chirp by UUID. |
| `POST` | `/api/chirps` | JWT | Creates a chirp for the authenticated user. |
| `DELETE` | `/api/chirps/{chirpID}` | JWT | Deletes an owned chirp. |
| `POST` | `/api/polka/webhooks` | API key | Processes a `user.upgraded` Polka event. |
| `POST` | `/admin/reset` | No | Deletes all users in `dev` mode only. |

Static files from the repository root are served below `/app/`.

## Example Workflow

Create a user:

```sh
curl -i -X POST http://localhost:8080/api/users \
  -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","password":"correct-horse-battery-staple"}'
```

Log in and save the returned `token` and `refresh_token`:

```sh
curl -i -X POST http://localhost:8080/api/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"ada@example.com","password":"correct-horse-battery-staple"}'
```

Create a chirp with the access token:

```sh
curl -i -X POST http://localhost:8080/api/chirps \
  -H 'Authorization: Bearer <access-token>' \
  -H 'Content-Type: application/json' \
  -d '{"body":"Hello from Chirpy"}'
```

Refresh an access token:

```sh
curl -i -X POST http://localhost:8080/api/refresh \
  -H 'Authorization: Bearer <refresh-token>'
```

Revoke a refresh token:

```sh
curl -i -X POST http://localhost:8080/api/revoke \
  -H 'Authorization: Bearer <refresh-token>'
```

## Polka Webhook

`POST /api/polka/webhooks` requires an `Authorization: ApiKey <POLKA_KEY>` header.

A successful upgrade event has this payload:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": "00000000-0000-0000-0000-000000000000"
  }
}
```

The referenced user is marked as a Chirpy Red subscriber.
