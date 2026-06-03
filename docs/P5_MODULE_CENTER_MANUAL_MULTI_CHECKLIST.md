# P5 Module Center Manual Multi Checklist

更新时间：2026-06-03

## 目标

P5.24 将 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 下的多模块人工验收入口收敛到模块中心页面的 checklist，让登录后验收不再只靠文档记忆。

本阶段不写数据库，不自动登录后台。

## 页面改动

模块中心的 registry checklist 新增：

- `多模块`：显示 `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1` 是否已经让 `Demo Notice` 进入 registry。
- `Demo Notice`：显示 `/demo/notice` 是否可作为打开入口。

模块中心阶段标识更新为 `P5.24`。

## Fixture 覆盖

`admin/src/views/dev_tools/module/registry-state.fixture.ts` 新增：

- `multiRegistryModules`
- `multiRegistryState`
- `multiChecklistRows`

覆盖 `Demo Article + Demo Notice` 同时存在时的页面 checklist 输出形状。

## 登录后人工验收

启动 API 时增加：

```bash
MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1
```

登录后台后打开：

```text
http://127.0.0.1:5173/module
```

预期关键字：

- `P5.24`
- `多模块`
- `MAKEADMIN_ENABLE_DEMO_NOTICE_MODULE=1`
- `Demo Notice`
- `/demo/notice`
- `未安装`

## No-DB 接入

新增：

```bash
scripts/check-module-center-manual-checklist.sh
```

并接入 `scripts/check-module-tools-no-db.sh`：

```text
Module tools: module center manual checklist contract
```

## 验收标准

- `scripts/check-module-center-manual-checklist.sh` 通过。
- `scripts/check-module-center-filter-contract.sh` 通过。
- `scripts/check-demo-notice-module.sh` 通过。
- `scripts/check-module-center-ui-contract.sh` 通过。
- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。

## 验收结果

- 已通过 `scripts/check-module-center-manual-checklist.sh`。
- 已通过 `scripts/check-module-center-filter-contract.sh`。
- 已通过 `scripts/check-demo-notice-module.sh`。
- 已通过 `scripts/check-module-center-ui-contract.sh`。
- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 已生成 `demo_notice` chunk。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.24 不做：

- 不自动登录后台。
- 不新增数据库 schema。
- 不修改 `.env`。
- 不处理 PTLM 业务模块。
