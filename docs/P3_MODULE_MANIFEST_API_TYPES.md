# P3 Module Manifest API Types

更新时间：2026-06-03

## 目标

P3.17 为管理端模块 manifest 预览、安装 apply、卸载 apply 补充前端 TypeScript 类型。

本阶段只收敛前端类型，不改变后端接口协议，不改变写入门禁，不创建或修改数据库 schema。

## 类型范围

`admin/src/api/tools/code.ts` 新增 manifest 相关类型：

- `ModuleManifestPreviewParams`
- `ModuleManifestInstallApplyParams`
- `ModuleManifestUninstallApplyParams`
- `ModuleManifestPreviewResult`
- `ModuleManifestApplyResult`
- `ModuleManifestPlanResult`
- `ModuleManifestApplySummaryResult`
- `ModuleManifestApplyCheckResult`
- `ModuleManifestApplySnapshotResult`

这些类型对应后端当前 JSON 字段，仅用于前端编译期约束。

## 页面收敛

`Manifest 预览` 弹窗现在使用 API 类型：

- `preview` 使用 `ModuleManifestPreviewResult`。
- `installResult` 和 `uninstallResult` 使用 `ModuleManifestApplyResult`。
- 安装、卸载请求参数分别使用 `ModuleManifestInstallApplyParams` 和 `ModuleManifestUninstallApplyParams`。
- 快照表格行使用本地 `SnapshotRow` 类型。

失败响应仍保留兜底对象，因为 request 拦截器在业务错误时会 reject 后端 `data`。

## 不在 P3.17 做

- 不修改请求 URL、method 或字段名。
- 不修改 request 封装。
- 不修改后端响应结构。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.17 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
