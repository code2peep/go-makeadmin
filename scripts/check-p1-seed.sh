#!/usr/bin/env bash
set -euo pipefail

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_PASSWORD="${MYSQL_PASSWORD-}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin_p1_check}"

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
    ma_admin
    ma_admin_profile
    ma_tenant
    ma_tenant_member
    ma_tenant_setting
    ma_role
    ma_admin_role
    ma_permission
    ma_role_permission
    ma_menu
    ma_menu_permission
    ma_org_unit
    ma_position
    ma_admin_org
    ma_data_scope
    ma_role_data_scope
    ma_login_log
    ma_audit_log
    ma_setting
    ma_dict_type
    ma_dict_item
    ma_file_category
    ma_file
    ma_codegen_table
    ma_codegen_column
)

failed=0
echo "==> Checking P1 ma_* tables"
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

admin_count="$(mysql_query "SELECT COUNT(*) FROM ma_admin WHERE username='admin' AND is_super=1 AND delete_time=0;")"
role_count="$(mysql_query "SELECT COUNT(*) FROM ma_role WHERE code='super_admin' AND tenant_id=0 AND delete_time=0;")"
permission_count="$(mysql_query "SELECT COUNT(*) FROM ma_permission;")"
role_permission_count="$(mysql_query "SELECT COUNT(*) FROM ma_role_permission WHERE role_id=1;")"
menu_count="$(mysql_query "SELECT COUNT(*) FROM ma_menu;")"
setting_count="$(mysql_query "SELECT COUNT(*) FROM ma_setting;")"
dict_type_count="$(mysql_query "SELECT COUNT(*) FROM ma_dict_type;")"
dict_item_count="$(mysql_query "SELECT COUNT(*) FROM ma_dict_item;")"
file_category_count="$(mysql_query "SELECT COUNT(*) FROM ma_file_category;")"

echo "==> Checking P1 seed rows"
if [ "$admin_count" -ne 1 ]; then
    echo "FAIL: admin seed is missing or duplicated."
    failed=1
else
    echo "OK: admin seed exists."
fi

if [ "$role_count" -ne 1 ]; then
    echo "FAIL: super_admin role seed is missing or duplicated."
    failed=1
else
    echo "OK: super_admin role seed exists."
fi

if [ "$permission_count" -lt 79 ]; then
    echo "FAIL: permission seed count is too low: $permission_count."
    failed=1
else
    echo "OK: permission seed count=$permission_count."
fi

if [ "$role_permission_count" -ne "$permission_count" ]; then
    echo "FAIL: role permission count ($role_permission_count) does not match permission count ($permission_count)."
    failed=1
else
    echo "OK: super_admin grants all permissions."
fi

if [ "$menu_count" -lt 22 ]; then
    echo "FAIL: menu seed count is too low: $menu_count."
    failed=1
else
    echo "OK: menu seed count=$menu_count."
fi

if [ "$setting_count" -lt 12 ]; then
    echo "FAIL: setting seed count is too low: $setting_count."
    failed=1
else
    echo "OK: setting seed count=$setting_count."
fi

if [ "$dict_type_count" -lt 4 ] || [ "$dict_item_count" -lt 14 ]; then
    echo "FAIL: dict seed count is incomplete: types=$dict_type_count items=$dict_item_count."
    failed=1
else
    echo "OK: dict seed counts types=$dict_type_count items=$dict_item_count."
fi

if [ "$file_category_count" -lt 2 ]; then
    echo "FAIL: file category seed count is too low: $file_category_count."
    failed=1
else
    echo "OK: file category seed count=$file_category_count."
fi

if [ "$failed" -ne 0 ]; then
    exit 1
fi

echo "==> check-p1-seed completed"
