#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"
INIT_P1_DROP="${INIT_P1_DROP:-0}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-123456}"

if [[ ! "$MYSQL_DATABASE" =~ ^[A-Za-z0-9_]+$ ]]; then
    echo "FAIL: MYSQL_DATABASE must contain only letters, numbers, and underscores."
    exit 1
fi

if ! command -v mysql >/dev/null 2>&1; then
    echo "FAIL: mysql client is not installed."
    exit 1
fi

if ! command -v go >/dev/null 2>&1; then
    echo "FAIL: go is not installed."
    exit 1
fi

password_len="${#ADMIN_PASSWORD}"
if [ "$password_len" -lt 6 ] || [ "$password_len" -gt 72 ]; then
    echo "FAIL: ADMIN_PASSWORD length must be 6-72 bytes."
    exit 1
fi

mysql_root() {
    MYSQL_PWD="$MYSQL_PASSWORD" mysql \
        --host="$MYSQL_HOST" \
        --port="$MYSQL_PORT" \
        --user="$MYSQL_USER" \
        --batch \
        --skip-column-names \
        --execute="$1"
}

mysql_db() {
    MYSQL_PWD="$MYSQL_PASSWORD" mysql \
        --host="$MYSQL_HOST" \
        --port="$MYSQL_PORT" \
        --user="$MYSQL_USER" \
        --database="$MYSQL_DATABASE" \
        "$@"
}

schema_exists="$(mysql_root "SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME='${MYSQL_DATABASE}';")"
if [ "$schema_exists" = "1" ]; then
    ma_table_count="$(mysql_root "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA='${MYSQL_DATABASE}' AND TABLE_NAME LIKE 'ma\\_%';")"
    if [ "$ma_table_count" -gt 0 ] && [ "$INIT_P1_DROP" != "1" ]; then
        echo "FAIL: database '$MYSQL_DATABASE' already has ma_* tables. Set INIT_P1_DROP=1 to recreate it."
        exit 1
    fi
fi

if [ "$INIT_P1_DROP" = "1" ]; then
    mysql_root "DROP DATABASE IF EXISTS \`${MYSQL_DATABASE}\`;"
fi
mysql_root "CREATE DATABASE IF NOT EXISTS \`${MYSQL_DATABASE}\` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

password_hash="$(cd "$ROOT_DIR/server" && MAKEADMIN_PASSWORD="$ADMIN_PASSWORD" go run ./cmd/makeadmin-password)"
seed_file="$(mktemp)"
trap 'rm -f "$seed_file"' EXIT

awk -v hash="$password_hash" '{
    gsub("INSTALL_TIME_PASSWORD_BCRYPT_REPLACE_ME", hash)
    print
}' "$ROOT_DIR/sql/p1.seed.sql" > "$seed_file"

mysql_db < "$ROOT_DIR/sql/p1.schema.sql"
mysql_db < "$seed_file"

MYSQL_HOST="$MYSQL_HOST" \
MYSQL_PORT="$MYSQL_PORT" \
MYSQL_USER="$MYSQL_USER" \
MYSQL_PASSWORD="$MYSQL_PASSWORD" \
MYSQL_DATABASE="$MYSQL_DATABASE" \
    "$ROOT_DIR/scripts/check-p1-seed.sh"

echo "==> init-p1-db completed for database '$MYSQL_DATABASE'"
