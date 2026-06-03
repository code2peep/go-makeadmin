# P5 Module Center Multi Filters

更新时间：2026-06-03

## 目标

P5.23 验证 `Demo Article` 与 `Demo Notice` 同时存在时，模块中心的状态筛选和统计不再依赖页面内联逻辑，避免多模块接入后 `全部`、`未安装`、`异常` 等筛选结果漂移。

本阶段不写数据库，不自动登录后台。

## 管理端改动

`admin/src/views/dev_tools/module/registry-state.ts` 新增：

```text
filterRegistryModules()
buildModuleStatusSummary()
isRegistryModuleFailed()
```

模块中心页面改为调用 helper：

- `filteredModules` 由 `filterRegistryModules()` 计算。
- `moduleStatusSummary` 由 `buildModuleStatusSummary()` 计算。
- 模块中心阶段标识更新为 `P5.23`。

## Fixture 覆盖

`admin/src/views/dev_tools/module/registry-state.fixture.ts` 新增双模块状态：

- `article`：`installed`
- `demo_notice`：`uninstalled`

编译期 fixture 覆盖：

- `multiStatusSummaryRows`
- `multiAllModules`
- `multiUninstalledModules`
- `multiFailedModules`

预期含义：

- 全部模块数量为 2。
- 未安装筛选命中 `Demo Notice`。
- 异常筛选在两个模块都通过 registry 校验时为空。

## No-DB 接入

新增：

```bash
scripts/check-module-center-filter-contract.sh
```

并接入 `scripts/check-module-tools-no-db.sh`：

```text
Module tools: module center filter contract
```

## 验收标准

- `scripts/check-module-center-filter-contract.sh` 通过。
- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## 验收结果

- 已通过 `scripts/check-module-center-filter-contract.sh`。
- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.23 不做：

- 不自动登录后台。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不处理 PTLM 业务模块。
