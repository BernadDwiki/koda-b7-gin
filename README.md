# KODA B7 — Backend (Gin)

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)

Backend service for the KODA B7 e-wallet exercise built with Go + Gin.

## Overview
- Lightweight REST API using Gin, PostgreSQL, Redis, and JWT for auth.
- Includes migrations, seeders, and OpenAPI/Swagger docs under `/swagger`.

## Technologies
- Go (modules)
- Gin (http framework)
- PostgreSQL (database)
- Redis (cache/session)
- JWT (authentication)

## Features
- Authentication: register, login, logout
- User profile: view and update (supports image upload)
- Wallet: balance, top-up, transfer
- Transactions: history, reporting
- API documentation: Swagger

## Requirements
- Go 1.20+ (or compatible)
- PostgreSQL 12+ (tested with 17)
- Redis
- Git

## Environment
Create a `.env` file in the project root with at least the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=ewallet
DB_SSL_MODE=disable
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=your_jwt_secret
PORT=8080
```

Adjust values to your environment. `DB_SSL_MODE` is used by the DB connection.

## Quickstart (local)
1. Clone the repository:

```bash
git clone <repo-url>
cd koda-b7-weekly10
```

2. Install dependencies:

```bash
go mod download
```

3. Apply migrations (project uses SQL files under `db/migrations`):

Use your preferred migration runner or run the SQL files against the DB directly. Example (psql):

```bash
psql "postgresql://$DB_USER:$DB_PASS@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSL_MODE" -f db/migrations/000001_create_users_table.up.sql
# repeat for other migration files in order
```

4. Run the server:

```bash
go run cmd/main.go
```

The server listens on `:8080` (or the value in `PORT`). API docs are available at `/swagger/index.html` when served.

## Docker (local)
Build and run with Docker using the repository root `Dockerfile`:

```bash
docker build -t koda-b7-backend .
docker run --env-file .env -p 8080:8080 koda-b7-backend
```

Or use `docker-compose.yml` at root (recommended for local dev with Postgres + Redis).

## CI / GHCR
This repository includes a GitHub Actions workflow `./.github/workflows/build_to_ghcr.yml` that builds and publishes an image to GHCR. Review and adjust the workflow for your repository (branch name, tags, and secrets).

## Main Routes (examples)
| Endpoint | Method | Description |
|---|---:|---|
| `/auth` | `POST` | Login |
| `/auth/new` | `POST` | Register |
| `/users/me` | `GET` | Get current user |
| `/profile` | `PUT` | Update profile (multipart/form-data for avatar) |
| `/wallet/top-up` | `POST` | Top up wallet |
| `/wallet/transfer` | `POST` | Create transfer (requires `pin`) |

Refer to the Swagger docs for a complete list of endpoints and request/response schemas.

## Seeders
Seeder SQL files are under `db/seeders/`. Run them against your DB to populate initial data.

## Contributing
- Fork the repo
- Create a feature branch
- Submit a pull request with a clear description and tests (if applicable)

## License
MIT License — see `LICENSE`.

## Contact
Maintainers: KODA B7 team
