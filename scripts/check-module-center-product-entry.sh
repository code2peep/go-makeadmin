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
STATUS_DOC="$ROOT/docs/P6_STATUS.md"
README="$ROOT/README.md"

assert_contains "$PAGE" "P6.1"
assert_contains "$PAGE" "module-market-overview"
assert_contains "$PAGE" "module-detail-panel"
assert_contains "$PAGE" "module-install-wizard"
assert_contains "$PAGE" "模块市场"
assert_contains "$PAGE" "模块详情"
assert_contains "$PAGE" "安装向导"
assert_contains "$HELPER" "buildModuleMarketRows"
assert_contains "$HELPER" "buildModuleInstallWizardRows"
assert_contains "$FIXTURE" "marketRows"
assert_contains "$FIXTURE" "selectedWizardRows"
assert_contains "$STATUS_DOC" "P6.1"
assert_contains "$README" "P6 状态"

echo "OK: module center product entry contract passed"
