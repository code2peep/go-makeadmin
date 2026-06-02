#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT_DIR/server"
MAKEADMIN_CODEGEN_FRONTEND_CHECK=1 \
GOCACHE="${GOCACHE:-/private/tmp/go-makeadmin-gocache}" \
go test ./generator -run TestGeneratedCrudFrontendCodeTypeChecks -count=1
