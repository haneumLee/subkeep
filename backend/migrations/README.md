# Database Migration Strategy for SubKeep

## Overview

SubKeep uses **golang-migrate** for database migrations with support for both up (forward) and down (rollback) migrations.

## Directory Structure

```
backend/
├── migrations/
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_subscriptions_table.up.sql
│   ├── 000002_create_subscriptions_table.down.sql
│   └── README.md
└── seeds/
    ├── dev/
    │   └── 001_sample_data.sql
    └── prod/
        └── 001_initial_data.sql
```

## Tools

### golang-migrate

```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir backend/migrations -seq create_users_table

# Run migrations
migrate -path backend/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback
migrate -path backend/migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" down 1

# Force version (if stuck)
migrate -path backend/migrations -database "..." force VERSION
```

## Migration Naming Convention

```
{version}_{description}.{up|down}.sql

Examples:
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_add_email_verification.up.sql
000002_add_email_verification.down.sql
```

## Best Practices

1. **Always create both up and down migrations**
2. **Test migrations thoroughly in development**
3. **Never modify existing migrations** - create new ones instead
4. **Use transactions** for data migrations
5. **Keep migrations small and focused**
6. **Backup before running in production**

## Integration with Code

### Option 1: Automatic (GORM AutoMigrate) - Development Only

```go
// DO NOT USE IN PRODUCTION
db.AutoMigrate(&User{}, &Subscription{})
```

### Option 2: golang-migrate in Code

```go
import (
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(dbURL string) error {
    m, err := migrate.New(
        "file://migrations",
        dbURL,
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}
```

## CI/CD Integration

Migrations run automatically in deployment pipeline:

```yaml
# .github/workflows/deploy.yml
- name: Run database migrations
  run: |
    migrate -path backend/migrations \
            -database "$DATABASE_URL" \
            up
```

## Seeding Data

```bash
# Development
psql -U user -d dbname -f backend/seeds/dev/001_sample_data.sql

# Production (initial data only)
psql -U user -d dbname -f backend/seeds/prod/001_initial_data.sql
```

## Emergency Rollback

```bash
# Rollback last migration
migrate -path backend/migrations -database "$DB_URL" down 1

# Rollback to specific version
migrate -path backend/migrations -database "$DB_URL" goto VERSION

# Force version (if migration is dirty)
migrate -path backend/migrations -database "$DB_URL" force VERSION
```

## Version Control

- **Keep migrations in version control**
- **Never delete migrations**
- **Document breaking changes**
- **Coordinate with team on migration order**
