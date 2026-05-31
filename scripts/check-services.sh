#!/usr/bin/env bash
set -euo pipefail

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"
REDIS_URL="${REDIS_URL:-redis://127.0.0.1:6379/0}"

failed=0

if [[ ! "$MYSQL_DATABASE" =~ ^[A-Za-z0-9_]+$ ]]; then
    echo "FAIL: MYSQL_DATABASE must contain only letters, numbers, and underscores."
    exit 1
fi

echo "==> Checking MySQL"
if command -v mysql >/dev/null 2>&1; then
    if MYSQL_PWD="$MYSQL_PASSWORD" mysql \
        --host="$MYSQL_HOST" \
        --port="$MYSQL_PORT" \
        --user="$MYSQL_USER" \
        --batch \
        --skip-column-names \
        --execute="SELECT SCHEMA_NAME FROM information_schema.SCHEMATA WHERE SCHEMA_NAME='${MYSQL_DATABASE}'" \
        | grep -qx "$MYSQL_DATABASE"; then
        echo "OK: MySQL is reachable and database '$MYSQL_DATABASE' exists."
    else
        echo "FAIL: MySQL is not reachable or database '$MYSQL_DATABASE' does not exist."
        failed=1
    fi
else
    echo "FAIL: mysql client is not installed; MySQL was not checked."
    failed=1
fi

echo "==> Checking Redis"
if command -v redis-cli >/dev/null 2>&1; then
    if redis-cli -u "$REDIS_URL" PING | grep -qx "PONG"; then
        echo "OK: Redis is reachable."
    else
        echo "FAIL: Redis is not reachable."
        failed=1
    fi
else
    echo "FAIL: redis-cli is not installed; Redis was not checked."
    failed=1
fi

if [ "$failed" -ne 0 ]; then
    echo "Service check failed. This script does not create databases, import SQL, or write Redis keys."
    exit 1
fi

echo "==> check-services completed"
