#!/bin/sh

# 1. Default DB URL
DB_URL=${POSTGRES_DSN:-"postgres://postgres:008008@postgres_container:5432/scheduling?sslmode=disable"}

# 2. Portable path detection (Works in sh, ash, dash, and zsh)
# This gets the directory where the script is located
SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
MIGRATIONS_PATH="$SCRIPT_DIR/../migrations"

# 3. Validation
if [ ! -d "$MIGRATIONS_PATH" ]; then
    echo "Error: Migration directory not found at $MIGRATIONS_PATH"
    exit 1
fi

case "$1" in
    up)
        echo "Running migrations up from $MIGRATIONS_PATH..."
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" up
        ;;
    down)
        echo "Rolling back last migration..."
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" down 1
        ;;
    force)
        echo "Forcing migration version $2..."
        migrate -path "$MIGRATIONS_PATH" -database "$DB_URL" force "$2"
        ;;
    *)
        echo "Usage: $0 {up|down|force version}"
        exit 1
esac