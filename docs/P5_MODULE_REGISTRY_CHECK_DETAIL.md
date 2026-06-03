# P5 Module Registry Check Detail

更新时间：2026-06-03

## 目标

P5.7 在模块中心展示 registry 与 manifest 的校验明细，让模块接入者可以直接看到每个 `manifestChecks` 检查项。

本阶段不新增后端接口，不新增数据库写入，不调整 manifest 校验规则。

## 管理端改动

模块中心 `校验` 列新增 `明细` 操作：

- 展示模块名、manifest 路径、整体校验状态和整体说明。
- 展示 `manifestChecks` 列表。
- 检查项状态映射为 `通过`、`异常`、`阻断` 或原始状态。
- 没有检查项时按钮不可点击，弹窗表格保留空态。

## 数据来源

明细全部来自 `GET /api/gen/moduleRegistry` 已返回的字段：

- `manifestStatus`
- `manifestMessage`
- `manifestChecks`

## 验收结果

- 已通过 `cd admin && npm run type-check`。
- 已通过 `GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh`。
- 全量验证里的前端 build 仍有 Rolldown 对 `@vueuse/core` pure annotation 的已知 warning，命令退出码为 0。
- 浏览器当前登录态已过期；页面人工复验需要重新登录后执行。

## 保留边界

P5.7 不做：

- 不新增 registry 写入接口。
- 不新增后端 API。
- 不从数据库读取模块 registry。
- 不创建或迁移业务 schema。
- 不处理 PTLM 业务模块。
