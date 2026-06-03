#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

expect_fail_no_db() {
    local output
    local status
    set +e
    output="$("$@" 2>&1)"
    status=$?
    set -e
    if [ "$status" -eq 0 ]; then
        echo "FAIL: expected command to fail before database access: $*"
        exit 1
    fi
    if [[ "$output" != *"no database access was attempted"* ]]; then
        echo "FAIL: expected no database access message from: $*"
        echo "$output"
        exit 1
    fi
}

expect_fail_no_db python3 scripts/module-codegen-plan.py --apply
expect_fail_no_db env MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 python3 scripts/module-codegen-plan.py --apply
expect_fail_no_db env MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 python3 scripts/module-codegen-plan.py \
    --apply \
    --confirm-module article
expect_fail_no_db env MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 python3 scripts/module-codegen-plan.py \
    --apply \
    --confirm-module article \
    --confirm-source-table ma_demo_article
expect_fail_no_db env MAKEADMIN_ALLOW_MODULE_CODEGEN_WRITE=1 python3 scripts/module-codegen-plan.py \
    --apply \
    --confirm-module article \
    --confirm-source-table wrong_table

echo "OK: module codegen apply boundary completed."
