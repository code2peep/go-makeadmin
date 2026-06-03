#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [ "${MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE:-}" != "1" ]; then
    echo "FAIL: module uninstall apply smoke requires MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE=1; no database access was attempted."
    exit 1
fi

cd "$ROOT/server"
: "${GOCACHE:=/private/tmp/go-makeadmin-gocache}"

if [ -z "${DATABASE_URL:-}" ]; then
    MYSQL_HOST="${MYSQL_HOST:-127.0.0.1}"
    MYSQL_PORT="${MYSQL_PORT:-3306}"
    MYSQL_USER="${MYSQL_USER:-root}"
    MYSQL_DATABASE="${MYSQL_DATABASE:-go_makeadmin}"
    MYSQL_AUTH="$MYSQL_USER"
    if [ -n "${MYSQL_PASSWORD:-}" ]; then
        MYSQL_AUTH="$MYSQL_USER:$MYSQL_PASSWORD"
    fi
    export DATABASE_URL="${MYSQL_AUTH}@tcp(${MYSQL_HOST}:${MYSQL_PORT})/${MYSQL_DATABASE}?charset=utf8mb4&parseTime=True&loc=Local"
fi

MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1 \
MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE=1 \
GOCACHE="$GOCACHE" \
go test ./generator/service/gen -run TestModuleManifestUninstallApplyLocalSmoke -count=1

echo "OK: module uninstall apply smoke completed."
