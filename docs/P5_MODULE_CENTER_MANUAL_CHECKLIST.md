# P5 Module Center Manual Checklist

更新时间：2026-06-03

## 目标

P5.17 将模块中心登录后人工验收项压缩到页面内的 registry checklist，减少反复翻 P5 文档才能判断当前模块清单状态的问题。

本阶段不新增后端接口，不新增测试依赖，不伪造登录态。

## 管理端改动

模块中心内置模块清单区域新增 `registry-manual-checklist`：

- `默认 Registry`：显示 `/api/gen/moduleRegistry` 当前读取状态。
- `Broken Fixture`：显示异常 fixture 是否随 API 返回。
- `异常筛选`：显示当前 registry 异常数量。
- `校验明细`：显示 manifestChecks 是否可打开。
- `Demo 入口`：显示 Demo Article 页面入口是否可打开。

模块中心阶段标识更新为 `P5.17`。

## Helper 改动

`admin/src/views/dev_tools/module/registry-state.ts` 新增：

```text
buildRegistryManualChecklistRows()
```

该 helper 只处理 registry 页面状态，不访问 API，不依赖 Vue 组件实例。

`admin/src/views/dev_tools/module/registry-state.fixture.ts` 已覆盖：

- 默认 registry checklist。
- broken fixture registry checklist。
- 空 registry checklist。

## 验收标准

- `cd admin && npm run type-check` 通过。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh` 通过。
- 登录后模块中心显示 `P5.17` 和 registry checklist。
- checklist 在默认 registry、broken fixture、空态和错误态下文案不撑破布局。

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前位于登录页，旧 token 已过期；本机未设置 `ADMIN_PASSWORD` 或 `P1_SMOKE_ADMIN_PASSWORD`，因此登录后页面截图验收需要你重新登录后执行。

## 保留边界

P5.17 不做：

- 不自动登录后台。
- 不修改 `.env` 或管理员密码。
- 不新增 registry 写入接口。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
