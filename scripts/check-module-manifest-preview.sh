#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOCACHE="${GOCACHE:-/private/tmp/go-makeadmin-gocache}"

(
    cd "$ROOT/server"
    GOCACHE="$GOCACHE" go test ./generator/service/gen -run TestPreviewModuleManifest -count=1
)

echo "OK: module manifest preview completed."
