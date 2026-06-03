# P5 Module Center UI Contract

更新时间：2026-06-03

## 目标

P5.18 为模块中心 registry UI 增加不依赖登录态的源码契约检查，确保默认 registry、broken fixture、异常筛选、校验明细和 Demo 入口这些关键验收文案不会在后续改页面时被漏删。

本阶段不做浏览器自动登录，不伪造 token，不修改管理员密码。

## 新增脚本

```bash
scripts/check-module-center-ui-contract.sh
```

脚本检查：

- 模块中心页面阶段标识为 `P5.21`。
- 页面包含 `registry-manual-checklist`。
- helper 保留 `buildRegistryManualChecklistRows()`。
- helper 保留 `默认 Registry`、`Broken Fixture`、`异常筛选`、`校验明细`、`Demo 入口`。
- fixture 覆盖 `buildRegistryManualChecklistRows()`。
- 本文档保留默认 registry 和 broken fixture 验收说明。

## No-DB 接入

`scripts/check-module-tools-no-db.sh` 已接入：

```text
Module tools: module center UI contract
```

因此 `./scripts/verify-no-db.sh` 会覆盖该检查。

## 登录后人工验收关键字

默认 registry 页面应能看到：

- `P5.21`
- `默认 Registry`
- `未开启`
- `Demo 入口`
- `/demo/article`

broken fixture 页面应能看到：

- `P5.21`
- `Broken Fixture`
- `已开启`
- `异常筛选`
- `校验明细`

## 验收标准

- `scripts/check-module-center-ui-contract.sh` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后模块中心可按上方关键字完成一次人工确认。

## 验收结果

- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.18 不做：

- 不自动登录后台。
- 不新增浏览器自动化依赖。
- 不修改 `.env` 或管理员密码。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
