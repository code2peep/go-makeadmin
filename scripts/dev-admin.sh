#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ADMIN_PORT="${ADMIN_PORT:-5173}"

echo "Starting admin dev server on http://127.0.0.1:$ADMIN_PORT"
cd "$ROOT_DIR/admin"
npm run dev -- --host 0.0.0.0 --port "$ADMIN_PORT" --strictPort "$@"
