#!/bin/bash
# ====================================
# Database Migration Management Script
# ====================================

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Configuration
DB_URL="${DATABASE_URL:-postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE:-disable}}"
MIGRATIONS_DIR="./migrations"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Functions
function print_help() {
    echo "Database Migration Management"
    echo ""
    echo "Usage: $0 [command] [options]"
    echo ""
    echo "Commands:"
    echo "  up              Apply all pending migrations"
    echo "  down [N]        Rollback last N migrations (default: 1)"
    echo "  goto VERSION    Migrate to specific version"
    echo "  force VERSION   Force set version (use with caution)"
    echo "  version         Show current migration version"
    echo "  create NAME     Create new migration files"
    echo "  seed [env]      Load seed data (dev|prod)"
    echo ""
    echo "Examples:"
    echo "  $0 up"
    echo "  $0 down 2"
    echo "  $0 create add_user_roles"
    echo "  $0 seed dev"
}

function check_migrate() {
    if ! command -v migrate &> /dev/null; then
        echo -e "${RED}Error: golang-migrate not found${NC}"
        echo "Install: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
        exit 1
    fi
}

function migrate_up() {
    echo -e "${BLUE}Applying migrations...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" up
    echo -e "${GREEN}✓ Migrations applied${NC}"
}

function migrate_down() {
    local steps=${1:-1}
    echo -e "${YELLOW}Rolling back $steps migration(s)...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" down "$steps"
    echo -e "${GREEN}✓ Rollback complete${NC}"
}

function migrate_goto() {
    local version=$1
    if [ -z "$version" ]; then
        echo -e "${RED}Error: Version required${NC}"
        exit 1
    fi
    echo -e "${BLUE}Migrating to version $version...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" goto "$version"
    echo -e "${GREEN}✓ Migration complete${NC}"
}

function migrate_force() {
    local version=$1
    if [ -z "$version" ]; then
        echo -e "${RED}Error: Version required${NC}"
        exit 1
    fi
    echo -e "${YELLOW}Forcing version to $version...${NC}"
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" force "$version"
    echo -e "${GREEN}✓ Version forced${NC}"
}

function show_version() {
    migrate -path "$MIGRATIONS_DIR" -database "$DB_URL" version
}

function create_migration() {
    local name=$1
    if [ -z "$name" ]; then
        echo -e "${RED}Error: Migration name required${NC}"
        exit 1
    fi
    migrate create -ext sql -dir "$MIGRATIONS_DIR" -seq "$name"
    echo -e "${GREEN}✓ Migration files created${NC}"
}

function load_seeds() {
    local env=${1:-dev}
    local seed_file="./seeds/$env/001_sample_data.sql"
    
    if [ ! -f "$seed_file" ]; then
        echo -e "${RED}Error: Seed file not found: $seed_file${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}Loading $env seed data...${NC}"
    psql "$DB_URL" -f "$seed_file"
    echo -e "${GREEN}✓ Seed data loaded${NC}"
}

# Main
check_migrate

case "$1" in
    up)
        migrate_up
        ;;
    down)
        migrate_down "$2"
        ;;
    goto)
        migrate_goto "$2"
        ;;
    force)
        migrate_force "$2"
        ;;
    version)
        show_version
        ;;
    create)
        create_migration "$2"
        ;;
    seed)
        load_seeds "$2"
        ;;
    *)
        print_help
        exit 1
        ;;
esac
