#!/usr/bin/env bash
set -euo pipefail

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"

if [[ ! "$MYSQL_DATABASE" =~ ^[A-Za-z0-9_]+$ ]]; then
    echo "FAIL: MYSQL_DATABASE must contain only letters, numbers, and underscores."
    exit 1
fi

if ! command -v mysql >/dev/null 2>&1; then
    echo "FAIL: mysql client is not installed."
    exit 1
fi

mysql_query() {
    MYSQL_PWD="$MYSQL_PASSWORD" mysql \
        --host="$MYSQL_HOST" \
        --port="$MYSQL_PORT" \
        --user="$MYSQL_USER" \
        --database="$MYSQL_DATABASE" \
        --batch \
        --skip-column-names \
        --execute="$1"
}

schema_exists="$(MYSQL_PWD="$MYSQL_PASSWORD" mysql \
    --host="$MYSQL_HOST" \
    --port="$MYSQL_PORT" \
    --user="$MYSQL_USER" \
    --batch \
    --skip-column-names \
    --execute="SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME='${MYSQL_DATABASE}';")"

if [ "$schema_exists" != "1" ]; then
    echo "FAIL: database '$MYSQL_DATABASE' does not exist."
    exit 1
fi

required_tables=(
    la_system_auth_admin
    la_system_auth_dept
    la_system_auth_menu
    la_system_auth_perm
    la_system_auth_role
    la_system_config
    la_system_log_login
)

failed=0
echo "==> Checking required blueprint tables"
for table in "${required_tables[@]}"; do
    exists="$(mysql_query "SELECT COUNT(*) FROM information_schema.TABLES WHERE TABLE_SCHEMA='${MYSQL_DATABASE}' AND TABLE_NAME='${table}';")"
    if [ "$exists" = "1" ]; then
        echo "OK: $table"
    else
        echo "FAIL: missing $table"
        failed=1
    fi
done

if [ "$failed" -ne 0 ]; then
    exit 1
fi

admin_count="$(mysql_query "SELECT COUNT(*) FROM la_system_auth_admin WHERE username='admin' AND is_delete=0;")"
menu_count="$(mysql_query "SELECT COUNT(*) FROM la_system_auth_menu;")"
config_count="$(mysql_query "SELECT COUNT(*) FROM la_system_config;")"

echo "==> Checking seed rows"
if [ "$admin_count" -lt 1 ]; then
    echo "FAIL: default admin seed is missing."
    failed=1
else
    echo "OK: default admin seed exists."
fi

if [ "$menu_count" -lt 1 ]; then
    echo "FAIL: menu seed is missing."
    failed=1
else
    echo "OK: menu seed count=$menu_count."
fi

if [ "$config_count" -lt 1 ]; then
    echo "FAIL: system config seed is missing."
    failed=1
else
    echo "OK: system config seed count=$config_count."
fi

if [ "$failed" -ne 0 ]; then
    exit 1
fi

echo "==> check-db-seed completed"
