# P5 Demo Notice Status Contract

更新时间：2026-06-03

## 目标

P5.22 将模块中心的运行时状态文案抽成纯 helper，并用 `Demo Notice` 覆盖 `runtimeRegistered=false` 的未注册状态，避免后续多模块回读时状态显示漂移。

本阶段不写数据库，不新增 `ma_demo_notice` 表，不注册后端运行时路由。

## 管理端改动

`admin/src/views/dev_tools/module/registry-state.ts` 新增：

```text
buildModuleRuntimeStatus()
```

输出契约：

| 输入 | 状态 | 说明 |
| --- | --- | --- |
| `runtimeRegistered=false` | `未注册` | 使用 `runtimeHint` 作为 detail |
| `runtimeRegistered=true` 且 env 未开启 | `未开启` | detail 为 `<ENV>=1` |
| 已注册且 env 开启 | `已开启` | detail 为 `<ENV>=1` 或 runtime hint |

`admin/src/views/dev_tools/module/index.vue` 不再内联运行时状态判断，改为调用 helper。

## Fixture 覆盖

`admin/src/views/dev_tools/module/registry-state.fixture.ts` 新增：

- `demoNoticeRuntimeStatus`
- `demoArticleRuntimeStatus`

其中 `demoNoticeRuntimeStatus` 覆盖 `runtimeRegistered=false` 的未注册状态。

## No-DB 接入

`scripts/check-demo-notice-module.sh` 已检查：

- 模块中心阶段标识为 `P5.22`。
- helper 包含 `buildModuleRuntimeStatus()`。
- helper 保留 `未注册` 文案。
- fixture 覆盖 `demoNoticeRuntimeStatus`。

该脚本已接入 `scripts/check-module-tools-no-db.sh`。

## 验收标准

- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## 验收结果

- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.22 不做：

- 不自动登录后台。
- 不新增 `ma_demo_notice` 表。
- 不注册 `/demo_notice/*` 后端运行时路由。
- 不修改 `.env`。
- 不处理 PTLM 业务模块。
