#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TMP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/go-makeadmin-module-codegen-plan.XXXXXX")"

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

python3 scripts/module-codegen-plan.py \
    --manifest "$TMP_DIR/manifest.json" \
    --tenant-id 0 \
    --format json >"$TMP_DIR/codegen-plan.json"

python3 - "$TMP_DIR/codegen-plan.json" <<'PY'
import json
import sys

plan = json.load(open(sys.argv[1]))
table = plan["makeadmin"]["table"]
columns = plan["makeadmin"]["columns"]
legacy = plan["legacy"]["genTable"]
if table["sourceTable"] != "ma_billing_invoice":
    raise SystemExit(f"unexpected source table: {table['sourceTable']}")
if table["templateType"] != "crud" or table["generateType"] != "zip":
    raise SystemExit("unexpected table generator defaults")
if legacy["genTpl"] != "crud" or legacy["genType"] != 0:
    raise SystemExit("unexpected legacy generator defaults")
if [column["columnName"] for column in columns] != ["id", "title", "status"]:
    raise SystemExit("unexpected generated columns")
title = next(column for column in columns if column["columnName"] == "title")
if title["queryType"] != "LIKE" or title["htmlType"] != "input":
    raise SystemExit("unexpected title column config")
PY

python3 - "$TMP_DIR/manifest.json" "$TMP_DIR/configured-manifest.json" <<'PY'
import json
import sys

manifest = json.load(open(sys.argv[1]))
manifest["codegen"] = {
    "columns": [
        {
            "columnName": "tenant_id",
            "goField": "tenantId",
            "htmlType": "input",
            "dictType": "",
        },
        {
            "columnName": "source_code",
            "goField": "sourceCode",
            "htmlType": "select",
            "dictType": "source_codes",
        },
    ]
}
json.dump(manifest, open(sys.argv[2], "w"), ensure_ascii=False, indent=2)
PY

python3 scripts/module-codegen-plan.py \
    --manifest "$TMP_DIR/configured-manifest.json" \
    --tenant-id 0 \
    --format json >"$TMP_DIR/configured-codegen-plan.json"

python3 - "$TMP_DIR/configured-codegen-plan.json" <<'PY'
import json
import sys

plan = json.load(open(sys.argv[1]))
columns = plan["makeadmin"]["columns"]
if [column["columnName"] for column in columns] != ["id", "tenant_id", "source_code"]:
    raise SystemExit("unexpected configured generated columns")
source = next(column for column in columns if column["columnName"] == "source_code")
if source["htmlType"] != "select" or source["dictType"] != "source_codes":
    raise SystemExit("unexpected configured source_code column config")
PY

echo "OK: module codegen plan completed."
