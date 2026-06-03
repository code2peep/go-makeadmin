#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

assert_contains() {
    local file="$1"
    local needle="$2"
    if ! grep -Fq "$needle" "$file"; then
        echo "FAIL: expected $file to contain: $needle"
        exit 1
    fi
}

PAGE="$ROOT/admin/src/views/workbench/index.vue"
SERVICE="$ROOT/server/admin/service/common/index.go"
DOC="$ROOT/docs/P6_STATUS.md"
README="$ROOT/README.md"

assert_contains "$PAGE" "P6.3 通用后台框架交付面"
assert_contains "$PAGE" "通用后台页面"
assert_contains "$PAGE" "AI 业务生成入口"
assert_contains "$PAGE" "模块中心"
assert_contains "$PAGE" "开发工具：manifest"
assert_contains "$SERVICE" "P6.3 通用后台框架交付面"
assert_contains "$SERVICE" "AI CRUD scaffold + codegen"
assert_contains "$SERVICE" "通用框架交付"
assert_contains "$DOC" "P6.3"
assert_contains "$README" "通用后台框架交付面"

echo "OK: framework workbench contract passed"
