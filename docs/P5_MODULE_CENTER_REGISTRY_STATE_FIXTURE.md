# P5 Module Center Registry State Fixture

更新时间：2026-06-03

## 目标

P5.15 为模块中心 registry 状态 helper 增加不依赖测试框架的 TypeScript 编译期 fixture，让 `vue-tsc` 覆盖默认、broken fixture、错误态和空态输入输出形状。

本阶段不新增前端测试依赖。

## 管理端改动

新增：

```text
admin/src/views/dev_tools/module/registry-state.fixture.ts
```

fixture 覆盖：

- 默认 registry module 输入。
- broken fixture registry module 输入。
- 空 registry 状态输入。
- 失败 registry 状态输入。
- `buildRegistryAcceptanceRows()` 输出形状。
- `countRegistryFailures()` 输出形状。
- `hasBrokenRegistryFixture()` 输出形状。
- `isRegistryEmptyState()` 输出形状。
- `registryTableEmptyTextFromState()` 输出形状。
- `registryErrorDetailText()` 输出形状。

模块中心阶段标识更新为 `P5.15`。

## 验证策略

fixture 位于 `admin/src/**/*`，会被 `admin/tsconfig.json` 的 include 规则纳入：

```bash
cd admin && npm run type-check
```

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后页面验收待重新登录后执行。

## 保留边界

P5.15 不做：

- 不新增 Vitest/Jest 等测试依赖。
- 不新增后端接口。
- 不伪造登录 token。
- 不修改 `.env` 或管理员密码。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
