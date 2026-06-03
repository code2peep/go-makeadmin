#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

assert_contains() {
    local file="$1"
    local needle="$2"
    if ! grep -Fq "$needle" "$file"; then
        echo "FAIL: expected $file to contain: $needle"
        exit 1
    fi
}

HELPER="$ROOT/admin/src/views/dev_tools/module/registry-state.ts"
FIXTURE="$ROOT/admin/src/views/dev_tools/module/registry-state.fixture.ts"
PAGE="$ROOT/admin/src/views/dev_tools/module/index.vue"
DOC="$ROOT/docs/P5_MODULE_CENTER_MANUAL_MULTI_CHECKLIST.md"

assert_contains "$PAGE" "P5.24"
assert_contains "$PAGE" "registry-manual-checklist"
assert_contains "$HELPER" "多模块"
assert_contains "$HELPER" "Demo Notice"
assert_contains "$HELPER" "MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1"
assert_contains "$HELPER" "/demo/notice"
assert_contains "$FIXTURE" "multiRegistryModules"
assert_contains "$FIXTURE" "multiChecklistRows"
assert_contains "$DOC" "MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1"
assert_contains "$DOC" "Demo Notice"
assert_contains "$DOC" "未安装"

echo "OK: module center manual checklist contract passed"
