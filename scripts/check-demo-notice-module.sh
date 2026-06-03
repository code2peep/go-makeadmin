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

echo "==> Demo Notice: manifest validation"
python3 "$ROOT/scripts/check-module-manifests.py" >/dev/null

echo "==> Demo Notice: service preview and registry"
(
    cd "$ROOT/server"
    MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1 \
        GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" \
        go test ./generator/service/gen \
        -run '^Test(ListModuleRegistryIncludesDemoNoticeWhenEnabled|DemoNoticeManifestUsesNoRuntimeGate|PreviewDemoNoticeManifestIncludesInstallPlan)$' \
        -count=1
)

echo "==> Demo Notice: route registry"
(
    cd "$ROOT/server"
    MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1 \
        GOPROXY="${GOPROXY:-https://goproxy.cn,direct}" \
        go test ./generator/routers/gen \
        -run '^TestListModuleRegistryRouteDemoNoticeResponse$' \
        -count=1
)

echo "==> Demo Notice: frontend entry contract"
assert_contains "$ROOT/admin/src/router/index.ts" "views/demo_notice"
assert_contains "$ROOT/admin/src/views/demo_notice/index.vue" "P5.21"
assert_contains "$ROOT/admin/src/views/demo_notice/index.vue" "Demo Notice"
assert_contains "$ROOT/admin/src/views/demo_notice/index.vue" "unregistered"
assert_contains "$ROOT/admin/src/views/dev_tools/module/index.vue" "P5.25"
assert_contains "$ROOT/admin/src/views/dev_tools/module/registry-state.ts" "buildModuleRuntimeStatus"
assert_contains "$ROOT/admin/src/views/dev_tools/module/registry-state.ts" "未注册"
assert_contains "$ROOT/admin/src/views/dev_tools/module/registry-state.fixture.ts" "demoNoticeRuntimeStatus"
assert_contains "$ROOT/examples/demo_notice/manifest.json" "\"routePath\": \"/dev_tools/demo/notice\""
assert_contains "$ROOT/examples/demo_notice/manifest.json" "\"component\": \"demo_notice/index\""

echo "OK: demo notice module contract passed"
