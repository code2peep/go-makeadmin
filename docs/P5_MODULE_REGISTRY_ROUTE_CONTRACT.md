# P5 Module Registry Route Contract

更新时间：2026-06-03

## 目标

P5.10 验证 `GET /api/gen/moduleRegistry` handler 的响应契约，确保前端模块中心依赖的 JSON 字段在路由层保持稳定。

由于当前本机登录态已过期，且没有可用的 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD` 环境变量，本阶段不做登录后浏览器截图验收。

## 后端改动

新增路由层测试：

```text
server/generator/routers/gen/module_registry_test.go
```

覆盖：

- 默认响应 `code=200`。
- 默认 registry 只返回 `Demo Article`。
- 默认 registry 的 `manifestStatus=passed`。
- 开启 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 后，响应包含 `Broken Manifest Fixture`。
- broken fixture 的 `manifestStatus=failed`。
- broken fixture 不影响 Demo Article 通过校验。

## Smoke 接入

`scripts/check-module-registry-smoke.sh` 已增加路由响应测试：

```bash
go test ./generator/routers/gen \
    -run '^TestListModuleRegistryRoute(DefaultResponse|BrokenFixtureResponse)$' \
    -count=1
```

因此 `./scripts/verify-no-db.sh` 会覆盖该路由契约。

## 管理端改动

模块中心阶段标识更新为 `P5.10`。

## 验收结果

- 已通过 `cd server && go test ./generator/routers/gen`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。

## 保留边界

P5.10 不做：

- 不绕过后台认证。
- 不伪造登录 token。
- 不修改 `.env` 或管理员密码。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
