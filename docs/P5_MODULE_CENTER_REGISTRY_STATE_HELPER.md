# P5 Module Center Registry State Helper

更新时间：2026-06-03

## 目标

P5.14 把模块中心 registry 状态计算从 `.vue` 页面中抽出为纯 TypeScript helper，为后续补前端单测或继续调整页面状态降低回归风险。

当前前端没有 Vitest/Jest 等单测框架，本阶段不新增测试依赖。

## 管理端改动

新增：

```text
admin/src/views/dev_tools/module/registry-state.ts
```

抽出的纯逻辑包括：

- registry 失败数量计算。
- broken fixture 是否存在。
- registry 空态判断。
- registry 错误详情文案。
- 表格 empty text 选择。
- 验收辅助状态条 rows 构造。

`admin/src/views/dev_tools/module/index.vue` 改为只负责请求、状态持有和渲染。

模块中心阶段标识更新为 `P5.14`。

## 验证策略

由于当前项目没有前端单测框架，P5.14 使用现有验证入口覆盖：

- `cd admin && npm run type-check`：验证 helper 类型契约。
- `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`：验证全量不触库链路。

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 登录后页面验收待重新登录后执行。

## 保留边界

P5.14 不做：

- 不新增前端测试依赖。
- 不新增后端接口。
- 不伪造登录 token。
- 不修改 `.env` 或管理员密码。
- 不新增数据库 schema。
- 不处理 PTLM 业务模块。
