#!/bin/bash

# Migration script for Go Backend Template
# Supports different environments and provides safety checks

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
MIGRATION_PATH="${MIGRATION_PATH:-./internal/data/migrations}"
ENVIRONMENT="${APP_SERVER_ENV:-development}"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to build database URL from environment variables
build_database_url() {
    local host="${APP_DATABASE_HOST:-localhost}"
    local port="${APP_DATABASE_PORT:-5432}"
    local user="${APP_DATABASE_USER:-postgres}"
    local password="${APP_DATABASE_PASSWORD:-dev_password}"
    local name="${APP_DATABASE_NAME:-golang_template_dev}"
    local sslmode="${APP_DATABASE_SSLMODE:-disable}"
    
    if [ "$ENVIRONMENT" = "production" ] && [ -z "$APP_DATABASE_PASSWORD" ]; then
        print_error "Database password is required for production environment"
        print_error "Set APP_DATABASE_PASSWORD environment variable"
        exit 1
    fi
    
    echo "postgres://${user}:${password}@${host}:${port}/${name}?sslmode=${sslmode}"
}

# Function to check if migrate CLI is installed
check_migrate_cli() {
    if ! command -v migrate &> /dev/null; then
        print_error "migrate CLI is not installed"
        print_error "Install it with: go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

# Function to check migration path exists
check_migration_path() {
    if [ ! -d "$MIGRATION_PATH" ]; then
        print_error "Migration path does not exist: $MIGRATION_PATH"
        exit 1
    fi
}

# Function to confirm production operations
confirm_production() {
    if [ "$ENVIRONMENT" = "production" ]; then
        print_warning "You are about to run migration on PRODUCTION environment!"
        print_warning "Database: $(build_database_url | sed 's/:.*@/:***@/')"
        read -p "Are you sure? (yes/no): " -r
        if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
            print_status "Migration cancelled"
            exit 0
        fi
    fi
}

# Function to show current migration status
show_status() {
    local db_url=$(build_database_url)
    print_status "Migration status for environment: $ENVIRONMENT"
    print_status "Database: $(echo $db_url | sed 's/:.*@/:***@/')"
    migrate -path "$MIGRATION_PATH" -database "$db_url" version
}

# Function to run migrations up
migrate_up() {
    local steps=${1:-""}
    local db_url=$(build_database_url)
    
    confirm_production
    
    print_status "Running migrations UP for environment: $ENVIRONMENT"
    if [ -n "$steps" ]; then
        print_status "Running $steps migration(s)"
        migrate -path "$MIGRATION_PATH" -database "$db_url" up "$steps"
    else
        print_status "Running all pending migrations"
        migrate -path "$MIGRATION_PATH" -database "$db_url" up
    fi
    print_status "Migration completed successfully"
}

# Function to run migrations down
migrate_down() {
    local steps=${1:-1}
    local db_url=$(build_database_url)
    
    confirm_production
    
    print_warning "Running migrations DOWN for environment: $ENVIRONMENT"
    print_warning "This will rollback $steps migration(s)"
    
    if [ "$ENVIRONMENT" = "production" ]; then
        read -p "Confirm rollback in production (yes/no): " -r
        if [[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]; then
            print_status "Migration cancelled"
            exit 0
        fi
    fi
    
    migrate -path "$MIGRATION_PATH" -database "$db_url" down "$steps"
    print_status "Migration rollback completed"
}

# Function to create new migration
create_migration() {
    local name="$1"
    if [ -z "$name" ]; then
        print_error "Migration name is required"
        print_error "Usage: $0 create <migration_name>"
        exit 1
    fi
    
    print_status "Creating new migration: $name"
    migrate create -ext sql -dir "$MIGRATION_PATH" -seq "$name"
    print_status "Migration files created in $MIGRATION_PATH"
}

# Function to force migration version
force_version() {
    local version="$1"
    if [ -z "$version" ]; then
        print_error "Version is required"
        print_error "Usage: $0 force <version>"
        exit 1
    fi
    
    local db_url=$(build_database_url)
    
    print_warning "Forcing migration version to $version for environment: $ENVIRONMENT"
    confirm_production
    
    migrate -path "$MIGRATION_PATH" -database "$db_url" force "$version"
    print_status "Migration version forced to $version"
}

# Function to show help
show_help() {
    echo "Migration script for Go Backend Template"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  status              Show current migration status"
    echo "  up [steps]          Run migrations up (all or specific number)"
    echo "  down [steps]        Run migrations down (default: 1)"
    echo "  create <name>       Create new migration"
    echo "  force <version>     Force migration version"
    echo "  help               Show this help"
    echo ""
    echo "Environment Variables:"
    echo "  APP_SERVER_ENV           Environment (development/production)"
    echo "  APP_DATABASE_HOST        Database host"
    echo "  APP_DATABASE_PORT        Database port"
    echo "  APP_DATABASE_USER        Database user"
    echo "  APP_DATABASE_PASSWORD    Database password"
    echo "  APP_DATABASE_NAME        Database name"
    echo "  APP_DATABASE_SSLMODE     SSL mode"
    echo "  MIGRATION_PATH           Path to migration files"
    echo ""
    echo "Examples:"
    echo "  $0 status                    # Show migration status"
    echo "  $0 up                        # Run all pending migrations"
    echo "  $0 up 1                      # Run one migration"
    echo "  $0 down                      # Rollback one migration"
    echo "  $0 down 2                    # Rollback two migrations"
    echo "  $0 create add_users_table    # Create new migration"
    echo ""
    echo "Production Example:"
    echo "  APP_SERVER_ENV=production \\"
    echo "  APP_DATABASE_HOST=prod-db.example.com \\"
    echo "  APP_DATABASE_PASSWORD=secret \\"
    echo "  $0 up"
}

# Main script logic
main() {
    check_migrate_cli
    check_migration_path
    
    case "${1:-help}" in
        "status")
            show_status
            ;;
        "up")
            migrate_up "$2"
            ;;
        "down")
            migrate_down "$2"
            ;;
        "create")
            create_migration "$2"
            ;;
        "force")
            force_version "$2"
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

main "$@"