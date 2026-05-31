#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SQL_FILE="${SQL_FILE:-$ROOT_DIR/sql/install.core.sql}"
MAKEADMIN_INIT_COMMAND="${MAKEADMIN_INIT_COMMAND:-$ROOT_DIR/scripts/init-local-core-db.sh}"
export SQL_FILE
export MAKEADMIN_INIT_COMMAND

exec "$ROOT_DIR/scripts/init-local-blueprint-db.sh" "$@"
