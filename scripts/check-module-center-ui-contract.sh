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

PAGE="$ROOT/admin/src/views/dev_tools/module/index.vue"
HELPER="$ROOT/admin/src/views/dev_tools/module/registry-state.ts"
FIXTURE="$ROOT/admin/src/views/dev_tools/module/registry-state.fixture.ts"
DOC="$ROOT/docs/P5_MODULE_CENTER_UI_CONTRACT.md"

assert_contains "$PAGE" "P5.18"
assert_contains "$PAGE" "registry-manual-checklist"
assert_contains "$HELPER" "buildRegistryManualChecklistRows"
assert_contains "$HELPER" "默认 Registry"
assert_contains "$HELPER" "Broken Fixture"
assert_contains "$HELPER" "异常筛选"
assert_contains "$HELPER" "校验明细"
assert_contains "$HELPER" "Demo 入口"
assert_contains "$FIXTURE" "buildRegistryManualChecklistRows"
assert_contains "$DOC" "默认 registry"
assert_contains "$DOC" "broken fixture"

echo "OK: module center UI contract passed"
