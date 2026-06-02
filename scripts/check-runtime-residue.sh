#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if ! command -v rg >/dev/null 2>&1; then
    echo "FAIL: rg is required for runtime residue guard."
    exit 1
fi

fail=0

check_no_match() {
    local label="$1"
    local pattern="$2"
    shift 2

    local output
    local status
    set +e
    output="$(rg -n --color=never --glob '*.go' --glob '*.tpl' "$pattern" "$@" 2>&1)"
    status=$?
    set -e
    if [ "$status" -eq 0 ]; then
        echo "FAIL: $label"
        echo "$output"
        fail=1
    elif [ "$status" -eq 1 ]; then
        echo "OK: $label"
    else
        echo "FAIL: $label scan failed"
        echo "$output"
        fail=1
    fi
}

cd "$ROOT_DIR"

echo "==> Checking P1 runtime residue"

runtime_paths=(
    server/admin/routers
    server/middleware
    server/makeadmin
    server/generator
)

check_no_match \
    "P1 runtime must not import old admin services or old system models" \
    'go-makeadmin/(admin/service/(system|setting|common)|model/(system|setting|common))' \
    "${runtime_paths[@]}"

check_no_match \
    "P1 runtime must not use old backstage Redis token keys" \
    'Backstage(Token|Manage|Roles|TokenSet)Key|backstage:' \
    server/admin/routers \
    server/middleware \
    server/makeadmin \
    server/generator

check_no_match \
    "P2 auth runtime must not use old makeadmin opaque token session keys" \
    'SessionToken(Key|Set)Prefix|makeadmin:token:' \
    "${runtime_paths[@]}"

check_no_match \
    "P2 tenant runtime must not hardcode GlobalTenantID in adapters or middleware" \
    'GlobalTenantID' \
    server/makeadmin/adapter \
    server/middleware

check_no_match \
    "P2 tenant-scoped settings, files and logs must not hardcode GlobalTenantID" \
    'GlobalTenantID' \
    server/makeadmin/repository/setting.go \
    server/makeadmin/repository/file.go \
    server/makeadmin/repository/log.go \
    server/makeadmin/service/setting.go \
    server/makeadmin/service/file.go \
    server/makeadmin/service/log.go

check_no_match \
    "P1 routes and middleware must not branch on adapter Available fallback" \
    '\.Available\(' \
    server/admin/routers \
    server/middleware

check_no_match \
    "P1 runtime must not reference la_* table names" \
    '\bla_[a-z0-9_]+\b' \
    "${runtime_paths[@]}"

check_no_match \
    "P1 runtime must not use legacy ConfigUtil setting fallback" \
    'ConfigUtil' \
    "${runtime_paths[@]}"

check_no_match \
    "config default table prefix must not fall back to la_" \
    'TablePrefix:[[:space:]]*"la_"|DB_TABLE_PREFIX[^"\n]*"la_"' \
    server/config

if [ "$fail" -ne 0 ]; then
    echo "FAIL: runtime residue guard found forbidden P1 references."
    exit 1
fi

echo "==> check-runtime-residue completed"
