#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

cd "$ROOT/server"

echo "==> Module registry smoke: default list"
unset MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE
GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" go test ./generator/service/gen \
    -run '^TestListModuleRegistryIncludesDemoArticle$' \
    -count=1

echo "==> Module registry smoke: broken fixture"
MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1 \
    GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" \
    go test ./generator/service/gen \
    -run '^TestListModuleRegistryIncludesBrokenFixtureWhenEnabled$' \
    -count=1

echo "OK: module registry default and broken fixture contracts passed"
