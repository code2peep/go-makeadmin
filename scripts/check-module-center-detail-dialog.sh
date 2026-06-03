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
STATUS_DOC="$ROOT/docs/P6_STATUS.md"
README="$ROOT/README.md"

assert_contains "$PAGE" "P6.2"
assert_contains "$PAGE" "module-detail-dialog"
assert_contains "$PAGE" "openModuleDetail"
assert_contains "$PAGE" "openSelectedModuleDetail"
assert_contains "$PAGE" "moduleDetailWizardRows"
assert_contains "$PAGE" "moduleDetailCheckRows"
assert_contains "$PAGE" "Manifest 校验"
assert_contains "$STATUS_DOC" "P6.2"
assert_contains "$README" "P6.2"

echo "OK: module center detail dialog contract passed"
