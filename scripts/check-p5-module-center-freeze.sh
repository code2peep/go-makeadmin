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

FINAL_DOC="$ROOT/docs/P5_FINAL_STATUS.md"
INDEX_DOC="$ROOT/docs/P5_MODULE_REGISTRY_DOC_INDEX.md"
STATUS_DOC="$ROOT/docs/P5_STATUS.md"
README="$ROOT/README.md"
PAGE="$ROOT/admin/src/views/dev_tools/module/index.vue"
HELPER="$ROOT/admin/src/views/dev_tools/module/registry-state.ts"
FIXTURE="$ROOT/admin/src/views/dev_tools/module/registry-state.fixture.ts"

assert_contains "$PAGE" "P5.25"
assert_contains "$README" "P5 最终状态"
assert_contains "$INDEX_DOC" "P5_FINAL_STATUS.md"
assert_contains "$STATUS_DOC" "P5.25"
assert_contains "$FINAL_DOC" "P5 Final Status"
assert_contains "$FINAL_DOC" "Demo Article"
assert_contains "$FINAL_DOC" "Demo Notice"
assert_contains "$FINAL_DOC" "scripts/check-module-registry-smoke.sh"
assert_contains "$FINAL_DOC" "scripts/check-module-center-ui-contract.sh"
assert_contains "$FINAL_DOC" "scripts/check-module-center-filter-contract.sh"
assert_contains "$FINAL_DOC" "scripts/check-module-center-manual-checklist.sh"
assert_contains "$FINAL_DOC" "scripts/check-demo-notice-module.sh"
assert_contains "$FINAL_DOC" "./scripts/verify-no-db.sh"
assert_contains "$HELPER" "buildRegistryManualChecklistRows"
assert_contains "$HELPER" "filterRegistryModules"
assert_contains "$HELPER" "buildModuleStatusSummary"
assert_contains "$FIXTURE" "multiChecklistRows"

echo "OK: P5 module center freeze contract passed"
