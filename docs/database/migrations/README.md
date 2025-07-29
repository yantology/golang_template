# Database Migrations Guide

This guide covers database migrations for different environments, with special focus on production deployment.

## 🚀 Quick Start

### Development Environment
```bash
# Start database
make db-up

# Run all migrations
make migrate-up

# Check status
make migrate-status
```

### Production Environment
```bash
# Set production environment variables
export APP_SERVER_ENV=production
export APP_DATABASE_HOST=your-prod-db.com
export APP_DATABASE_USER=your-user
export APP_DATABASE_PASSWORD=your-secure-password
export APP_DATABASE_NAME=your-prod-db

# Run migrations
make migrate-prod ARGS="up"
```

## 📋 Migration Commands

### Basic Commands (Makefile)
```bash
# Environment-based (uses DATABASE_URL or defaults to development)
make migrate-up                    # Run all pending migrations
make migrate-down                  # Rollback one migration
make migrate-status                # Check current status
make migrate-create NAME=add_table # Create new migration

# Advanced commands
make migrate-up-one                # Run one migration up
make migrate-down-one              # Run one migration down
make migrate-force VERSION=1       # Force migration version
```

### Advanced Commands (Helper Script)
```bash
# Using the migration script directly
./scripts/migrate.sh status                    # Show migration status
./scripts/migrate.sh up                        # Run all migrations
./scripts/migrate.sh up 2                      # Run 2 migrations
./scripts/migrate.sh down                      # Rollback 1 migration
./scripts/migrate.sh down 3                    # Rollback 3 migrations
./scripts/migrate.sh create add_users_table    # Create new migration
./scripts/migrate.sh force 5                   # Force version to 5

# Using make with script
make migrate ARGS="status"                     # Check status
make migrate ARGS="up"                         # Run migrations
make migrate ARGS="create add_products_table"  # Create migration
```

## 🏭 Production Migration Strategy

### 1. Pre-Production Testing
```bash
# Test on staging environment
APP_SERVER_ENV=staging \
APP_DATABASE_HOST=staging-db.example.com \
APP_DATABASE_PASSWORD=staging-password \
./scripts/migrate.sh up
```

### 2. Production Deployment
```bash
# Set production environment
export APP_SERVER_ENV=production
export APP_DATABASE_HOST=prod-db.example.com
export APP_DATABASE_USER=app_user
export APP_DATABASE_PASSWORD=your-secure-password
export APP_DATABASE_NAME=your_app_prod
export APP_DATABASE_SSLMODE=require

# Check current status
./scripts/migrate.sh status

# Run migrations with confirmation
./scripts/migrate.sh up
```

### 3. Production Safety Features
- **Environment confirmation**: Script asks for confirmation in production
- **Password validation**: Requires password for production environment
- **SSL enforcement**: Uses `sslmode=require` for production
- **Rollback protection**: Extra confirmation for down migrations

## 🔧 Environment Configuration

### Environment Variables
```bash
# Required for all environments
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_NAME=your_db

# Required for production
APP_DATABASE_PASSWORD=secure-password
APP_SERVER_ENV=production

# Optional
APP_DATABASE_SSLMODE=disable     # development: disable, production: require
MIGRATION_PATH=./internal/data/migrations
```

### Configuration Examples

#### Development (.env)
```bash
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USER=postgres
APP_DATABASE_PASSWORD=dev_password
APP_DATABASE_NAME=golang_template_dev
APP_DATABASE_SSLMODE=disable
```

#### Production (.env.production)
```bash
APP_DATABASE_HOST=prod-db.example.com
APP_DATABASE_PORT=5432
APP_DATABASE_USER=app_user
APP_DATABASE_PASSWORD=your-secure-password
APP_DATABASE_NAME=your_app_prod
APP_DATABASE_SSLMODE=require
APP_SERVER_ENV=production
```

## 📝 Creating Migrations

### 1. Create Migration Files
```bash
# Using make
make migrate-create NAME=add_users_table

# Using script
./scripts/migrate.sh create add_users_table

# Manual (with migrate CLI)
migrate create -ext sql -dir ./internal/data/migrations -seq add_users_table
```

### 2. Migration File Structure
```
internal/data/migrations/
├── 000001_create_users_table.up.sql    # Forward migration
├── 000001_create_users_table.down.sql  # Rollback migration
├── 000002_add_products_table.up.sql
└── 000002_add_products_table.down.sql
```

### 3. Migration Best Practices

#### ✅ Good Migration (000001_create_users_table.up.sql)
```sql
-- Create users table with proper constraints
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    email_verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);

-- Add trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at 
    BEFORE UPDATE ON users 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();
```

#### ✅ Good Rollback (000001_create_users_table.down.sql)
```sql
-- Drop in reverse order
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP INDEX IF EXISTS idx_users_active;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
```

### 4. Migration Guidelines
- **Always use `IF EXISTS` / `IF NOT EXISTS`** for idempotency
- **Create proper indexes** for query performance
- **Use UUIDs for primary keys** for distributed systems
- **Add timestamps** (created_at, updated_at) to all tables
- **Test rollbacks** - ensure down migrations work
- **Use transactions** for complex migrations
- **Avoid data migrations** in schema migrations

## 🚨 Production Migration Checklist

### Before Migration
- [ ] Test migration on staging environment
- [ ] Backup production database
- [ ] Schedule maintenance window (if needed)
- [ ] Notify team of migration
- [ ] Verify application compatibility

### During Migration
- [ ] Check current migration status
- [ ] Run migration with confirmation
- [ ] Monitor database performance
- [ ] Check application health after migration

### After Migration
- [ ] Verify migration status
- [ ] Test critical application features
- [ ] Monitor error logs
- [ ] Confirm application performance

### Rollback Plan
- [ ] Know which migration to rollback to
- [ ] Test rollback on staging first
- [ ] Have database backup ready
- [ ] Plan application deployment rollback

## 🛠️ Troubleshooting

### Common Issues

#### Migration Stuck
```bash
# Check current version
./scripts/migrate.sh status

# Force to specific version (use carefully!)
./scripts/migrate.sh force 5
```

#### Connection Issues
```bash
# Test database connection
psql "postgres://user:password@host:port/database?sslmode=require"

# Check environment variables
env | grep APP_DATABASE
```

#### Migration Conflicts
```bash
# Resolve conflicts by forcing version
./scripts/migrate.sh force <last_known_good_version>

# Then run migrations normally
./scripts/migrate.sh up
```

### Production Issues

#### SSL Certificate Issues
```bash
# For production, ensure SSL is configured
export APP_DATABASE_SSLMODE=require

# For development/testing, disable SSL
export APP_DATABASE_SSLMODE=disable
```

#### Permission Issues
```bash
# Ensure database user has proper permissions
GRANT CREATE, ALTER, DROP ON DATABASE your_db TO your_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO your_user;
```

## 📚 Additional Resources

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migration Best Practices](https://www.prisma.io/blog/database-migration-best-practices)