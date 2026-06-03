#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> Runtime residue guard"
"$ROOT_DIR/scripts/check-runtime-residue.sh"

echo "==> Framework workbench contract"
"$ROOT_DIR/scripts/check-framework-workbench-contract.sh"

echo "==> Module tools no-db guard"
"$ROOT_DIR/scripts/check-module-tools-no-db.sh"

echo "==> Backend: go test ./..."
(
    cd "$ROOT_DIR/server"
    GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" go test ./...
)

echo "==> Frontend: npm run type-check"
(
    cd "$ROOT_DIR/admin"
    npm run type-check
)

echo "==> Frontend: npm run build"
(
    cd "$ROOT_DIR/admin"
    npm run build
)

echo "==> Frontend: npm audit"
(
    cd "$ROOT_DIR/admin"
    npm audit --audit-level=moderate
)

echo "==> verify-no-db completed"
