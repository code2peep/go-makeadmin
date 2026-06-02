#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST=""
TMP_DIR=""

while [ "$#" -gt 0 ]; do
    case "$1" in
        --manifest)
            MANIFEST="${2:-}"
            shift 2
            ;;
        *)
            echo "FAIL: unknown argument: $1"
            exit 1
            ;;
    esac
done

if [ -z "$MANIFEST" ]; then
    TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/go-makeadmin-module-codegen.XXXXXX")"
fi

cleanup() {
    if [ -n "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}
trap cleanup EXIT

cd "$ROOT"

if [ -z "$MANIFEST" ]; then
    MANIFEST="$TMP_DIR/manifest.json"
    python3 scripts/module-scaffold.py \
        --module billing_invoice \
        --entity BillingInvoice \
        --table ma_billing_invoice \
        --requires-schema \
        --print-manifest >"$MANIFEST"
else
    MANIFEST="$(cd "$ROOT" && realpath "$MANIFEST")"
fi

python3 -m json.tool "$MANIFEST" >/dev/null

cd "$ROOT/server"
MAKEADMIN_CODEGEN_MANIFEST="$MANIFEST" \
GOCACHE="${GOCACHE:-/private/tmp/go-makeadmin-gocache}" \
go test ./generator -run TestGeneratedCrudCodeMatchesModuleManifest -count=1

echo "OK: module scaffold output matches codegen templates."
