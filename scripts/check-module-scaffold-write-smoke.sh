#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ "${MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE:-}" != "1" ]; then
    echo "FAIL: module scaffold write smoke requires MAKEADMIN_ALLOW_MODULE_SCAFFOLD_WRITE=1; no files were written and no database access was attempted."
    exit 1
fi

STAMP="$(date +%Y%m%d%H%M%S)"
MODULE="billing_invoice_${STAMP}"
ENTITY="BillingInvoice${STAMP}"
TABLE="ma_billing_invoice_${STAMP}"
EXAMPLES_ROOT=".cache/module-scaffold-smoke/${STAMP}/examples"
MANIFEST="${EXAMPLES_ROOT}/${MODULE}/manifest.json"
README="${EXAMPLES_ROOT}/${MODULE}/README.md"

cd "$ROOT"

python3 scripts/module-scaffold.py \
    --module "$MODULE" \
    --entity "$ENTITY" \
    --table "$TABLE" \
    --requires-schema \
    --examples-root "$EXAMPLES_ROOT" >/tmp/go-makeadmin-module-scaffold-write-smoke.out

test -f "$MANIFEST"
test -f "$README"

python3 -m json.tool "$MANIFEST" >/dev/null
python3 scripts/module-install-plan.py --manifest "$MANIFEST" --tenant-id 0 --role-id 1 >/dev/null
python3 scripts/module-uninstall-plan.py --manifest "$MANIFEST" >/dev/null
scripts/check-module-codegen.sh --manifest "$MANIFEST" >/dev/null

if ! grep -q "$MANIFEST" "$README"; then
    echo "FAIL: README does not reference generated manifest path"
    exit 1
fi

echo "OK: module scaffold write smoke completed at ${EXAMPLES_ROOT}/${MODULE}"
