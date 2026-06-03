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

echo "==> Module tools: python syntax"
python3 -m py_compile \
    scripts/check-module-manifests.py \
    scripts/module-scaffold.py \
    scripts/module-codegen-plan.py \
    scripts/module-registry-plan.py \
    scripts/module-role-grant-plan.py \
    scripts/module-install-plan.py \
    scripts/module-uninstall-plan.py

echo "==> Module tools: shell syntax"
bash -n scripts/check-module-lifecycle-smoke.sh scripts/check-module-codegen.sh scripts/check-module-codegen-plan.sh scripts/check-module-codegen-apply-boundary.sh scripts/check-module-codegen-apply-smoke.sh scripts/check-module-codegen-readback-smoke.sh scripts/check-module-install-plan-preview.sh scripts/check-module-install-apply-boundary.sh scripts/check-module-manifest-preview.sh scripts/check-module-scaffold-write-smoke.sh

echo "==> Module tools: manifest validation"
python3 scripts/check-module-manifests.py >/dev/null

echo "==> Module tools: scaffold dry-run"
python3 scripts/module-scaffold.py \
    --module billing_invoice \
    --entity BillingInvoice \
    --table ma_billing_invoice \
    --requires-schema \
    --dry-run >/dev/null

echo "==> Module tools: scaffold codegen link"
scripts/check-module-codegen.sh >/dev/null

echo "==> Module tools: scaffold codegen plan"
scripts/check-module-codegen-plan.sh >/dev/null

echo "==> Module tools: manifest preview"
scripts/check-module-manifest-preview.sh >/dev/null

echo "==> Module tools: install plan preview"
scripts/check-module-install-plan-preview.sh >/dev/null

echo "==> Module tools: install apply boundary"
scripts/check-module-install-apply-boundary.sh >/dev/null

echo "==> Module tools: codegen apply boundary"
scripts/check-module-codegen-apply-boundary.sh >/dev/null

echo "==> Module tools: dry-run previews"
python3 scripts/module-codegen-plan.py --manifest examples/demo/manifest.json --format json >/dev/null
python3 scripts/module-registry-plan.py --manifest examples/demo/manifest.json >/dev/null
python3 scripts/module-role-grant-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1 >/dev/null
python3 scripts/module-install-plan.py --manifest examples/demo/manifest.json --tenant-id 0 --role-id 1 >/dev/null
python3 scripts/module-uninstall-plan.py --manifest examples/demo/manifest.json >/dev/null

echo "==> Module tools: write gates"
expect_fail_no_db python3 scripts/module-registry-plan.py --apply
expect_fail_no_db python3 scripts/module-codegen-plan.py --apply
expect_fail_no_db python3 scripts/module-install-plan.py --apply
expect_fail_no_db python3 scripts/module-uninstall-plan.py --apply
expect_fail_no_db scripts/check-module-lifecycle-smoke.sh

echo "==> Module tools: filesystem write gates"
expect_fail_no_db scripts/check-module-scaffold-write-smoke.sh

echo "==> Module tools: database write smoke gates"
expect_fail_no_db scripts/check-module-codegen-apply-smoke.sh
expect_fail_no_db scripts/check-module-codegen-readback-smoke.sh

echo "==> check-module-tools-no-db completed"
