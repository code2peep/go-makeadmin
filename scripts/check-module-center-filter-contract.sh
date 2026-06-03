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
DOC="$ROOT/docs/P5_MODULE_CENTER_MULTI_FILTERS.md"

assert_contains "$HELPER" "filterRegistryModules"
assert_contains "$HELPER" "buildModuleStatusSummary"
assert_contains "$HELPER" "isRegistryModuleFailed"
assert_contains "$FIXTURE" "multiStatusModules"
assert_contains "$FIXTURE" "multiUninstalledModules"
assert_contains "$FIXTURE" "multiFailedModules"
assert_contains "$FIXTURE" "demo_notice"
assert_contains "$PAGE" "P5.24"
assert_contains "$PAGE" "filterRegistryModules"
assert_contains "$DOC" "Demo Article"
assert_contains "$DOC" "Demo Notice"

echo "OK: module center filter contract passed"
