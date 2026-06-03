#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST="examples/demo/manifest.json"
MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
MYSQL_PORT="${MYSQL_PORT:-3306}"
MYSQL_USER="${MYSQL_USER:-root}"
MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"

if [ "${MAKEADMIN_ALLOW_DEMO_MODULE_VISIBLE_WRITE:-}" != "1" ]; then
    echo "FAIL: demo module visible smoke requires MAKEADMIN_ALLOW_DEMO_MODULE_VISIBLE_WRITE=1; no database access was attempted."
    exit 1
fi

if ! command -v mysql >/dev/null 2>&1; then
    echo "FAIL: mysql client is required."
    exit 1
fi

query() {
    local sql="$1"
    if [ -n "${MYSQL_PASSWORD:-}" ]; then
        MYSQL_PWD="$MYSQL_PASSWORD" mysql \
            --host "$MYSQL_HOST" \
            --port "$MYSQL_PORT" \
            --user "$MYSQL_USER" \
            --database "$MYSQL_DATABASE" \
            --batch \
            --raw \
            --skip-column-names \
            --execute "$sql"
        return
    fi
    mysql \
        --host "$MYSQL_HOST" \
        --port "$MYSQL_PORT" \
        --user "$MYSQL_USER" \
        --database "$MYSQL_DATABASE" \
        --batch \
        --raw \
        --skip-column-names \
        --execute "$sql"
}

cd "$ROOT_DIR"

python3 scripts/check-module-manifests.py >/dev/null

manifest_flags="$(python3 - <<'PY'
import json
from pathlib import Path

manifest = json.loads(Path("examples/demo/manifest.json").read_text())
print(int(manifest["menu"]["visible"] is True))
print(int(manifest["runtimeRegistered"] is True))
print(manifest["menu"]["routePath"])
print(manifest["menu"]["component"])
PY
)"

if [ "$manifest_flags" != $'1\n1\n/dev_tools/demo/article\narticle/index' ]; then
    echo "FAIL: demo manifest must be visible, runtime registered, and point to /dev_tools/demo/article -> article/index."
    exit 1
fi

if [ ! -f "admin/src/api/article.ts" ] \
    || [ ! -f "admin/src/views/article/index.vue" ] \
    || [ ! -f "admin/src/views/article/edit.vue" ]; then
    echo "FAIL: demo article frontend files are missing."
    exit 1
fi

MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1 \
python3 scripts/module-uninstall-plan.py \
    --manifest "$MANIFEST" \
    --confirm-module article \
    --confirm-delete \
    --apply >/dev/null

MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1 \
python3 scripts/module-install-plan.py \
    --manifest "$MANIFEST" \
    --tenant-id 0 \
    --role-id 1 \
    --confirm-module article \
    --confirm-role-id 1 \
    --apply >/dev/null

counts="$(query "
SELECT
    (SELECT COUNT(*) FROM ma_permission WHERE code IN ('article:list','article:detail','article:add','article:edit','article:del')),
    (SELECT COUNT(*) FROM ma_menu WHERE route_name = 'demo.article' AND route_path = '/dev_tools/demo/article' AND component = 'article/index' AND is_visible = 1 AND status = 1 AND delete_time = 0),
    (
        SELECT COUNT(*) FROM ma_menu_permission AS mp
        INNER JOIN ma_menu AS m ON m.id = mp.menu_id
        INNER JOIN ma_permission AS p ON p.id = mp.permission_id
        WHERE m.route_name = 'demo.article'
          AND p.code = 'article:list'
    ),
    (
        SELECT COUNT(*) FROM ma_role_permission AS rp
        INNER JOIN ma_permission AS p ON p.id = rp.permission_id
        WHERE rp.tenant_id = 0
          AND rp.role_id = 1
          AND p.code IN ('article:list','article:detail','article:add','article:edit','article:del')
    );
")"

if [ "$counts" != $'5\t1\t1\t5' ]; then
    echo "FAIL: unexpected demo module install counts: $counts"
    exit 1
fi

echo "OK: demo module visible install completed."
