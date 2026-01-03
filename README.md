# Go Backend Template

A production-ready Go backend template with GraphQL, PostgreSQL, and dependency injection.

## Features

- **GraphQL API** - gqlgen with directives support
- **PostgreSQL** - sqlc for type-safe queries, goose for migrations
- **Dependency Injection** - uber-go/fx
- **Configuration** - viper with YAML + environment variable overrides
- **Authentication** - JWT/JWK validation (Firebase compatible)

## Prerequisites

- Go 1.21+
- PostgreSQL
- Make

## Quick Start

```bash
# Configure database
cp application.yaml application.local.yaml
# Edit application.local.yaml with your database credentials

# Run the server
make run-server
```

Server starts at `http://localhost:8888/query`

## Commands

```bash
make run-server    # Run the server
make sqlc          # Generate database code
make gqlgen        # Generate GraphQL code
```

## Project Structure

```
cmd/server/         # Application entry point
internal/
  bootstrap/        # Dependency injection setup
  graph/            # GraphQL schema and resolvers
  storage/pg/       # Database queries and migrations
```

See [CLAUDE.md](./CLAUDE.md) for detailed development guidelines.
