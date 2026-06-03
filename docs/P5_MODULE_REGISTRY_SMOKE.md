# P5 Module Registry Smoke

更新时间：2026-06-03

## 目标

P5.9 把模块 registry 的后端数据契约固化成 CLI smoke，让默认清单和 broken fixture 清单都可以在不登录后台的情况下验证。

本阶段不新增数据库写入，不新增 schema，不新增后端 API。

## 新增脚本

```bash
scripts/check-module-registry-smoke.sh
```

脚本验证两组契约：

- 默认清单：未设置 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE` 时，只返回 `Demo Article`，且 manifest 校验通过。
- 异常 fixture：设置 `MAKEADMIN_ENABLE_BROKEN_MODULE_REGISTRY_FIXTURE=1` 时，返回 `Broken Manifest Fixture`，且异常项校验失败但不影响 Demo Article。

该脚本已接入：

```bash
scripts/check-module-tools-no-db.sh
```

因此也会被以下全量验证覆盖：

```bash
./scripts/verify-no-db.sh
```

## 管理端改动

模块中心阶段标识更新为 `P5.9`。

## 验收结果

- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。

## 保留边界

P5.9 不做：

- 不依赖后台登录态做 registry 契约验证。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不默认启用 broken fixture。
- 不处理 PTLM 业务模块。
