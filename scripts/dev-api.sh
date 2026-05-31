#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

read_env_value() {
    local key="$1"
    local file="$ROOT_DIR/server/.env"
    if [ ! -f "$file" ]; then
        return 0
    fi
    awk -F= -v key="$key" '$1 == key {
        value = substr($0, index($0, "=") + 1)
        sub(/\r$/, "", value)
        print value
    }' "$file" | tail -n 1
}

SERVER_HOST="${SERVER_HOST:-$(read_env_value SERVER_HOST)}"
SERVER_PORT="${SERVER_PORT:-$(read_env_value SERVER_PORT)}"
SERVER_HOST="${SERVER_HOST:-127.0.0.1}"
SERVER_PORT="${SERVER_PORT:-8000}"

if command -v lsof >/dev/null 2>&1 && lsof -nP -iTCP:"$SERVER_PORT" -sTCP:LISTEN >/dev/null 2>&1; then
    echo "FAIL: TCP port $SERVER_PORT is already in use."
    lsof -nP -iTCP:"$SERVER_PORT" -sTCP:LISTEN
    echo "Change SERVER_PORT in server/.env and VITE_API_PROXY_TARGET in admin/.env.development."
    exit 1
fi

echo "Starting API server on http://$SERVER_HOST:$SERVER_PORT"
echo "This requires MySQL and Redis to be ready. Run ./scripts/check-services.sh first."
cd "$ROOT_DIR/server"
go run . "$@"
