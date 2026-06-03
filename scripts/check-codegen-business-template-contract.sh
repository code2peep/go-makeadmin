#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/go-makeadmin-codegen-contract.XXXXXX")"

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cd "$ROOT"

require_text() {
    local file="$1"
    local text="$2"
    if ! grep -Fq "$text" "$file"; then
        echo "FAIL: expected $file to contain: $text"
        exit 1
    fi
}

require_text server/generator/templates/vue/index.vue.tpl "resetPage"
require_text server/generator/templates/vue/index.vue.tpl "resetParams"
require_text server/generator/templates/vue/index.vue.tpl "新增"
require_text server/generator/templates/vue/index.vue.tpl "handleEdit"
require_text server/generator/templates/vue/index.vue.tpl "handleDelete"
require_text server/generator/templates/vue/index.vue.tpl "pagination"
require_text server/generator/templates/vue/index.vue.tpl "dict-value"
require_text server/generator/templates/vue/index.vue.tpl "useDictData"

require_text server/generator/templates/vue/edit.vue.tpl "Popup"
require_text server/generator/templates/vue/edit.vue.tpl "formRules"
require_text server/generator/templates/vue/edit.vue.tpl "handleSubmit"
require_text server/generator/templates/vue/edit.vue.tpl "getDetail"
require_text server/generator/templates/vue/edit.vue.tpl "defineExpose"
require_text server/generator/templates/vue/edit.vue.tpl "dictData"

require_text server/generator/templates/vue/api.ts.tpl "Lists"
require_text server/generator/templates/vue/api.ts.tpl "Detail"
require_text server/generator/templates/vue/api.ts.tpl "Add"
require_text server/generator/templates/vue/api.ts.tpl "Edit"
require_text server/generator/templates/vue/api.ts.tpl "Delete"

require_text server/generator/tpl_test.go "codegen_template_smoke"
require_text server/generator/tpl_test.go "common_status"
require_text docs/P6_STATUS.md "P6.4"

python3 scripts/module-codegen-plan.py \
    --manifest examples/demo/manifest.json \
    --tenant-id 0 \
    --format json >"$TMP_DIR/codegen-plan.json"

python3 - "$TMP_DIR/codegen-plan.json" <<'PY'
import json
import sys

plan = json.load(open(sys.argv[1]))
columns = plan["makeadmin"]["columns"]
status = next((column for column in columns if column["columnName"] == "status"), None)
if status is None:
    raise SystemExit("missing default status column")
if status["htmlType"] != "radio" or status["dictType"] != "common_status":
    raise SystemExit(f"unexpected default status column: {status}")
legacy = next((column for column in plan["legacy"]["genTableColumns"] if column["columnName"] == "status"), None)
if legacy is None:
    raise SystemExit("missing legacy status column")
if legacy["htmlType"] != "radio" or legacy["dictType"] != "common_status":
    raise SystemExit(f"unexpected legacy status column: {legacy}")
PY

echo "OK: codegen business template contract completed."
