#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ "${MAKEADMIN_ALLOW_MODULE_CODEGEN_SMOKE_WRITE:-}" != "1" ]; then
    echo "FAIL: module codegen apply smoke requires MAKEADMIN_ALLOW_MODULE_CODEGEN_SMOKE_WRITE=1; no database access was attempted."
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

query() {
    "${MYSQL[@]}" --execute "$1"
}

live_table_count_sql() {
    cat <<'SQL'
SELECT COUNT(*)
FROM ma_codegen_table
WHERE tenant_id = 0 AND table_name = 'ma_demo_article' AND delete_time = 0;
SQL
}

detail_sql() {
    cat <<'SQL'
SELECT
    COUNT(DISTINCT t.id),
    COUNT(c.id),
    COALESCE(GROUP_CONCAT(c.column_name ORDER BY c.sort SEPARATOR ','), '')
FROM ma_codegen_table AS t
LEFT JOIN ma_codegen_column AS c ON c.table_id = t.id
WHERE t.tenant_id = 0 AND t.table_name = 'ma_demo_article' AND t.delete_time = 0;
SQL
}

insert_stale_column_sql() {
    cat <<'SQL'
INSERT INTO ma_codegen_column
(table_id, column_name, column_comment, column_type, column_length, go_type, go_field, json_field,
 is_pk, is_increment, is_required, is_insert, is_edit, is_list, is_query, query_type, html_type,
 dict_type, sort, create_time, update_time)
SELECT id, 'stale_column', 'Stale Column', 'varchar', 64, 'string', 'staleColumn', 'staleColumn',
       0, 0, 0, 1, 1, 1, 0, '=', 'input', '', 99, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()
FROM ma_codegen_table
WHERE tenant_id = 0 AND table_name = 'ma_demo_article' AND delete_time = 0
LIMIT 1
ON DUPLICATE KEY UPDATE update_time = UNIX_TIMESTAMP();
SQL
}

cleanup_sql() {
    cat <<'SQL'
DELETE c FROM ma_codegen_column AS c
INNER JOIN ma_codegen_table AS t ON t.id = c.table_id
WHERE t.tenant_id = 0 AND t.table_name = 'ma_demo_article' AND t.delete_time = 0;

UPDATE ma_codegen_table
SET delete_time = UNIX_TIMESTAMP(), update_time = UNIX_TIMESTAMP()
WHERE tenant_id = 0 AND table_name = 'ma_demo_article' AND delete_time = 0;
SQL
}

start_count="$(query "$(live_table_count_sql)")"
if [ "$start_count" != "0" ]; then
    echo "FAIL: live ma_demo_article codegen rows already exist before smoke: $start_count"
    exit 1
fi

cleanup() {
    query "$(cleanup_sql)" >/dev/null
}
trap cleanup EXIT

cd "$ROOT"

MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 \
python3 scripts/module-codegen-plan.py \
    --manifest examples/demo/manifest.json \
    --tenant-id 0 \
    --confirm-module article \
    --confirm-source-table ma_demo_article \
    --confirm-sync-columns \
    --mysql-host "$MYSQL_HOST" \
    --mysql-port "$MYSQL_PORT" \
    --mysql-user "$MYSQL_USER" \
    --mysql-database "$MYSQL_DATABASE" \
    --apply

after_first="$(query "$(detail_sql)")"
echo "after_first=$after_first"
if [ "$after_first" != $'1\t3\tid,title,status' ]; then
    echo "FAIL: unexpected first apply counts"
    exit 1
fi

query "$(insert_stale_column_sql)" >/dev/null
after_stale="$(query "$(detail_sql)")"
echo "after_stale=$after_stale"
if [ "$after_stale" != $'1\t4\tid,title,status,stale_column' ]; then
    echo "FAIL: unexpected stale column counts"
    exit 1
fi

MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 \
python3 scripts/module-codegen-plan.py \
    --manifest examples/demo/manifest.json \
    --tenant-id 0 \
    --confirm-module article \
    --confirm-source-table ma_demo_article \
    --confirm-sync-columns \
    --mysql-host "$MYSQL_HOST" \
    --mysql-port "$MYSQL_PORT" \
    --mysql-user "$MYSQL_USER" \
    --mysql-database "$MYSQL_DATABASE" \
    --apply

after_second="$(query "$(detail_sql)")"
echo "after_second=$after_second"
if [ "$after_second" != $'1\t3\tid,title,status' ]; then
    echo "FAIL: unexpected second apply counts"
    exit 1
fi

cleanup
trap - EXIT

after_cleanup="$(query "$(live_table_count_sql)")"
if [ "$after_cleanup" != "0" ]; then
    echo "FAIL: cleanup left live codegen rows: $after_cleanup"
    exit 1
fi

echo "OK: module codegen apply smoke completed."
