#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "==> Backend menu tree contract"
(
    cd "$ROOT_DIR/server"
    GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" go test ./util ./makeadmin/adapter -run 'TestListToTree|TestRouteMenuMapsBuildsNestedChildren'
)

echo "==> Frontend catalogue redirect contract"
rg -q "createCatalogueRedirect" "$ROOT_DIR/admin/src/router/index.ts"
rg -q "findFirstChildRoutePath" "$ROOT_DIR/admin/src/router/index.ts"
rg -q "routeRecord.redirect" "$ROOT_DIR/admin/src/router/index.ts"
rg -q "buildFullRoutePath" "$ROOT_DIR/admin/src/router/index.ts"
rg -q "filterAsyncRoutes\\(route.children, false, fullPath\\)" "$ROOT_DIR/admin/src/router/index.ts"
rg -q "joinRoutePaths\\(fullPath, childPath\\)" "$ROOT_DIR/admin/src/router/index.ts"

echo "==> menu tree contract passed"
