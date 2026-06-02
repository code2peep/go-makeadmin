#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/go-makeadmin-module-codegen.XXXXXX")"

cleanup() {
    rm -rf "$TMP_DIR"
}
trap cleanup EXIT

cd "$ROOT"

python3 scripts/module-scaffold.py \
    --module billing_invoice \
    --entity BillingInvoice \
    --table ma_billing_invoice \
    --requires-schema \
    --print-manifest >"$TMP_DIR/manifest.json"

python3 -m json.tool "$TMP_DIR/manifest.json" >/dev/null

cd "$ROOT/server"
MAKEADMIN_CODEGEN_MANIFEST="$TMP_DIR/manifest.json" \
GOCACHE="${GOCACHE:-/private/tmp/go-makeadmin-gocache}" \
go test ./generator -run TestGeneratedCrudCodeMatchesModuleManifest -count=1

echo "OK: module scaffold output matches codegen templates."
