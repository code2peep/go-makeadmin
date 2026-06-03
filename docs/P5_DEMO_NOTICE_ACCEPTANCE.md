# P5 Demo Notice Acceptance

更新时间：2026-06-03

## 目标

P5.21 为第二个只读示例模块 `Demo Notice` 增加安装计划和页面入口验收，确认它不仅能出现在 registry 中，也能被 manifest preview、install plan、管理端路由入口和未注册运行时提示覆盖。

本阶段不写数据库，不新增 `ma_demo_notice` 表，不注册后端运行时路由。

## 新增脚本

```bash
scripts/check-demo-notice-module.sh
```

脚本覆盖：

- `examples/demo_notice/manifest.json` 通过 manifest 校验。
- `PreviewModuleManifest()` 能生成 `Demo Notice` 的安装计划。
- 安装计划包含 `demo.notice`、`demo_notice:list`、`demo_notice:detail` 和卸载 SQL。
- registry 在 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 下返回 `Demo Notice`。
- route handler 在同一 env 下返回 `Demo Notice`。
- 管理端动态路由包含 `demo_notice` 页面。
- `admin/src/views/demo_notice/index.vue` 显示 `P5.21`、`Demo Notice` 和 `unregistered`。
- 模块中心阶段标识更新为 `P5.21`。

## No-DB 接入

`scripts/check-module-tools-no-db.sh` 已接入：

```text
Module tools: demo notice module contract
```

因此 `./scripts/verify-no-db.sh` 会覆盖 Demo Notice 的 no-db 验收。

## 本地人工入口

启动 API 时增加：

```bash
MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1
```

登录后台后打开模块中心：

```text
http://127.0.0.1:5173/module
```

预期能看到 `Demo Notice`、`/demo/notice`、`未注册` 或未注册运行时相关提示。

## 验收标准

- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-registry-smoke.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## 验收结果

- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-registry-smoke.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.21 不做：

- 不自动登录后台。
- 不新增 `ma_demo_notice` 表。
- 不注册 `/demo_notice/*` 后端运行时路由。
- 不修改 `.env`。
- 不处理 PTLM 业务模块。
