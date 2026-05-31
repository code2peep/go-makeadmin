#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

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
