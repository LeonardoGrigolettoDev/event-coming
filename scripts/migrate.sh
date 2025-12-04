#!/bin/bash

# Migration script for Event-Coming database

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo -e "${RED}ERROR: DATABASE_URL environment variable is not set${NC}"
    echo "Example: export DATABASE_URL='postgresql://postgres:postgres@localhost:5432/event_coming?sslmode=disable'"
    exit 1
fi

# Check if migrate is installed
if ! command -v migrate &> /dev/null; then
    echo -e "${RED}ERROR: migrate tool is not installed${NC}"
    echo "Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    exit 1
fi

# Parse command line arguments
COMMAND=${1:-"up"}

case $COMMAND in
    up)
        echo -e "${GREEN}Running migrations UP...${NC}"
        migrate -path migrations -database "$DATABASE_URL" up
        echo -e "${GREEN}Migrations completed successfully!${NC}"
        ;;
    down)
        echo -e "${YELLOW}Rolling back migrations DOWN...${NC}"
        migrate -path migrations -database "$DATABASE_URL" down
        echo -e "${GREEN}Rollback completed successfully!${NC}"
        ;;
    force)
        if [ -z "$2" ]; then
            echo -e "${RED}ERROR: Version number required for force command${NC}"
            echo "Usage: $0 force <version>"
            exit 1
        fi
        echo -e "${YELLOW}Forcing version to $2...${NC}"
        migrate -path migrations -database "$DATABASE_URL" force $2
        echo -e "${GREEN}Version forced successfully!${NC}"
        ;;
    version)
        echo -e "${GREEN}Current migration version:${NC}"
        migrate -path migrations -database "$DATABASE_URL" version
        ;;
    *)
        echo "Usage: $0 {up|down|force <version>|version}"
        exit 1
        ;;
esac
