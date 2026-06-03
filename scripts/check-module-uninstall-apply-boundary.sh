#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT/server"

: "${GOCACHE:=/private/tmp/go-makeadmin-gocache}"

GOCACHE="$GOCACHE" go test ./generator/service/gen -run TestModuleManifestUninstallApplyGate -count=1

echo "OK: module uninstall apply boundary completed."
