#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ "${MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE:-}" != "1" ]; then
    echo "FAIL: module codegen readback smoke requires MAKEADMIN_ALLOW_MODULE_CODEGEN_READBACK_WRITE=1; no database access was attempted."
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
    echo "FAIL: live ma_demo_article codegen rows already exist before readback smoke: $start_count"
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

(
    cd server
    MAKEADMIN_CODEGEN_READBACK_SMOKE=1 \
    MYSQL_HOST="$MYSQL_HOST" \
    MYSQL_PORT="$MYSQL_PORT" \
    MYSQL_USER="$MYSQL_USER" \
    MYSQL_DATABASE="$MYSQL_DATABASE" \
    MYSQL_PASSWORD="${MYSQL_PASSWORD:-}" \
    go test ./generator/service/gen -run TestCodegenConfigReadbackAndTemplateGenerationSmoke -count=1
)

cleanup
trap - EXIT

after_cleanup="$(query "$(live_table_count_sql)")"
if [ "$after_cleanup" != "0" ]; then
    echo "FAIL: cleanup left live codegen rows: $after_cleanup"
    exit 1
fi

echo "OK: module codegen readback smoke completed."
