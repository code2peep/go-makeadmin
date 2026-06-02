#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ "${MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE:-}" != "1" ]; then
    echo "FAIL: module lifecycle smoke requires MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1; no database access was attempted."
    exit 1
fi

MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"
if [ -n "${MYSQL_PASSWORD:-}" ]; then
    export MYSQL_PWD="$MYSQL_PASSWORD"
fi

MYSQL=(
    mysql
    --host "$MYSQL_HOST"
    --port "$MYSQL_PORT"
    --user "$MYSQL_USER"
    --database "$MYSQL_DATABASE"
    --batch
    --raw
    --skip-column-names
)

CODES="'article:list','article:detail','article:add','article:edit','article:del'"

query() {
    "${MYSQL[@]}" --execute "$1"
}

total_count_sql() {
    cat <<SQL
SELECT
    (SELECT COUNT(*) FROM ma_permission WHERE code IN ($CODES)) +
    (SELECT COUNT(*) FROM ma_menu WHERE route_name = 'demo.article') +
    (
        SELECT COUNT(*) FROM ma_menu_permission mp
        LEFT JOIN ma_menu m ON m.id = mp.menu_id
        LEFT JOIN ma_permission p ON p.id = mp.permission_id
        WHERE m.route_name = 'demo.article' OR p.code IN ($CODES)
    ) +
    (
        SELECT COUNT(*) FROM ma_role_permission rp
        LEFT JOIN ma_permission p ON p.id = rp.permission_id
        WHERE p.code IN ($CODES)
    );
SQL
}

detail_counts_sql() {
    cat <<SQL
SELECT
    (SELECT COUNT(*) FROM ma_permission WHERE code IN ($CODES)),
    (SELECT COUNT(*) FROM ma_menu WHERE route_name = 'demo.article'),
    (
        SELECT COUNT(*) FROM ma_menu_permission mp
        LEFT JOIN ma_menu m ON m.id = mp.menu_id
        LEFT JOIN ma_permission p ON p.id = mp.permission_id
        WHERE m.route_name = 'demo.article' AND p.code = 'article:list'
    ),
    (
        SELECT COUNT(*) FROM ma_role_permission rp
        LEFT JOIN ma_permission p ON p.id = rp.permission_id
        WHERE p.code IN ($CODES)
    );
SQL
}

cleanup_sql() {
    cat <<SQL
DELETE rp FROM ma_role_permission rp
LEFT JOIN ma_permission p ON p.id = rp.permission_id
WHERE p.code IN ($CODES);
DELETE mp FROM ma_menu_permission mp
LEFT JOIN ma_menu m ON m.id = mp.menu_id
LEFT JOIN ma_permission p ON p.id = mp.permission_id
WHERE m.route_name = 'demo.article' OR p.code IN ($CODES);
DELETE FROM ma_menu WHERE route_name = 'demo.article';
DELETE FROM ma_permission WHERE code IN ($CODES);
SQL
}

start_count="$(query "$(total_count_sql)")"
if [ "$start_count" != "0" ]; then
    echo "FAIL: demo article rows already exist before lifecycle smoke: $start_count"
    exit 1
fi

cleanup() {
    query "$(cleanup_sql)" >/dev/null
}
trap cleanup EXIT

cd "$ROOT"

MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 \
python3 scripts/module-install-plan.py \
    --manifest examples/demo/manifest.json \
    --tenant-id 0 \
    --role-id 1 \
    --confirm-module article \
    --confirm-role-id 1 \
    --mysql-host "$MYSQL_HOST" \
    --mysql-port "$MYSQL_PORT" \
    --mysql-user "$MYSQL_USER" \
    --mysql-database "$MYSQL_DATABASE" \
    --apply

installed_counts="$(query "$(detail_counts_sql)")"
echo "after_install=$installed_counts"
if [ "$installed_counts" != $'5\t1\t1\t5' ]; then
    echo "FAIL: unexpected install counts"
    exit 1
fi

MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 \
python3 scripts/module-uninstall-plan.py \
    --manifest examples/demo/manifest.json \
    --confirm-module article \
    --confirm-delete \
    --mysql-host "$MYSQL_HOST" \
    --mysql-port "$MYSQL_PORT" \
    --mysql-user "$MYSQL_USER" \
    --mysql-database "$MYSQL_DATABASE" \
    --apply

after_uninstall="$(query "$(total_count_sql)")"
echo "after_uninstall_total=$after_uninstall"
if [ "$after_uninstall" != "0" ]; then
    echo "FAIL: uninstall left rows: $after_uninstall"
    exit 1
fi

MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 \
python3 scripts/module-uninstall-plan.py \
    --manifest examples/demo/manifest.json \
    --confirm-module article \
    --confirm-delete \
    --mysql-host "$MYSQL_HOST" \
    --mysql-port "$MYSQL_PORT" \
    --mysql-user "$MYSQL_USER" \
    --mysql-database "$MYSQL_DATABASE" \
    --apply

after_second_uninstall="$(query "$(total_count_sql)")"
if [ "$after_second_uninstall" != "0" ]; then
    echo "FAIL: second uninstall left rows: $after_second_uninstall"
    exit 1
fi

trap - EXIT
echo "OK: module lifecycle smoke completed."
