#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"
SQL_FILE="${SQL_FILE:-$ROOT_DIR/sql/install.sql}"
INIT_COMMAND="${MAKEADMIN_INIT_COMMAND:-$0}"
APPLY=0

case "${1:-}" in
    "")
        ;;
    "--dry-run")
        ;;
    "--apply")
        APPLY=1
        ;;
    *)
        echo "Usage: $0 [--dry-run|--apply]"
        exit 1
        ;;
esac

if [[ ! "$MYSQL_DATABASE" =~ ^[A-Za-z0-9_]+$ ]]; then
    echo "FAIL: MYSQL_DATABASE must contain only letters, numbers, and underscores."
    exit 1
fi

if [ ! -f "$SQL_FILE" ]; then
    echo "FAIL: SQL file not found: $SQL_FILE"
    exit 1
fi

if [ "$APPLY" -eq 0 ]; then
    echo "DRY RUN: no database changes will be made."
    echo "Target database: $MYSQL_DATABASE"
    echo "SQL file: $SQL_FILE"
    echo
    echo "This plan would:"
    echo "1. Create database '$MYSQL_DATABASE' if it does not exist."
    echo "2. Refuse to import when '$MYSQL_DATABASE' is non-empty unless MAKEADMIN_ALLOW_NONEMPTY_DB=1."
    echo "3. Import $SQL_FILE, which creates P0 la_* tables."
    echo
    echo "To apply after explicit approval:"
    echo "MAKEADMIN_ALLOW_SCHEMA_WRITE=1 $INIT_COMMAND --apply"
    exit 0
fi

if [ "${MAKEADMIN_ALLOW_SCHEMA_WRITE:-0}" != "1" ]; then
    echo "FAIL: schema writes are disabled."
    echo "Set MAKEADMIN_ALLOW_SCHEMA_WRITE=1 and pass --apply only after explicit approval."
    exit 1
fi

if ! command -v mysql >/dev/null 2>&1; then
    echo "FAIL: mysql client is not installed."
    exit 1
fi

mysql_exec() {
    MYSQL_PWD="$MYSQL_PASSWORD" mysql \
        --host="$MYSQL_HOST" \
        --port="$MYSQL_PORT" \
        --user="$MYSQL_USER" \
        --batch \
        --skip-column-names \
        "$@"
}

echo "==> Creating database '$MYSQL_DATABASE' when missing"
mysql_exec --execute="CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE} CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

table_count="$(mysql_exec --execute="SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA='${MYSQL_DATABASE}';")"
if [ "$table_count" != "0" ] && [ "${MAKEADMIN_ALLOW_NONEMPTY_DB:-0}" != "1" ]; then
    echo "FAIL: database '$MYSQL_DATABASE' already has $table_count tables."
    echo "Refusing to import because sql/install.sql contains DROP TABLE statements."
    echo "Use a disposable empty database, or set MAKEADMIN_ALLOW_NONEMPTY_DB=1 after explicit approval."
    exit 1
fi

echo "==> Importing SQL into '$MYSQL_DATABASE'"
MYSQL_PWD="$MYSQL_PASSWORD" mysql \
    --host="$MYSQL_HOST" \
    --port="$MYSQL_PORT" \
    --user="$MYSQL_USER" \
    "$MYSQL_DATABASE" < "$SQL_FILE"

echo "==> database initialization completed"
