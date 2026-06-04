#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

COMPONENT="$ROOT_DIR/admin/src/components/material/index.vue"
HOOK="$ROOT_DIR/admin/src/components/material/hook.ts"
SERVICE="$ROOT_DIR/server/admin/service/common/album.go"
DOC="$ROOT_DIR/docs/P6_STATUS.md"

rg -Fq "id: -1" "$HOOK"
rg -Fq "未分组" "$HOOK"
rg -Fq "listReq.Cid >= 0" "$SERVICE"
rg -Fq "material-left__header" "$COMPONENT"
rg -Fq "el-empty" "$COMPONENT"
rg -Fq "uploadButtonText" "$COMPONENT"
rg -Fq "cateId.value > 0 ? cateId.value : 0" "$COMPONENT"
rg -q "暂无.*上传后会显示在这里" "$COMPONENT"
rg -Fq "未分组暂无" "$COMPONENT"
rg -Fq "emptyDescription" "$COMPONENT"
rg -Fq "P6.8" "$DOC"

echo "OK: material empty state contract passed"
