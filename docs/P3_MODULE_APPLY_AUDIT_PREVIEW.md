# P3 Module Apply Audit Preview

更新时间：2026-06-03

## 目标

P3.22 为管理端模块 manifest apply 结果增加审计事件 dry-run 预览。

本阶段只在前端基于现有 apply result 生成 JSON 预览，不调用后端，不新增接口，不写数据库。

## 前端构造器

`admin/src/api/tools/code.ts` 新增：

```ts
buildModuleManifestApplyAuditPreview(result, options)
```

输入：

- `ModuleManifestApplyResult`
- 可选 `eventId`
- 可选 `actor`
- 可选 `requestedAt`
- 可选 `completedAt`

输出：

- `ModuleManifestApplyAuditEventResult`

## 映射规则

- `source`、`manifest`、`summary` 来自 apply result。
- `scope.tenantId` 优先来自 apply result，其次来自 plan。
- `scope.roleId` 优先来自 apply result，其次来自 plan。
- `scope.databaseScope` 来自 summary。
- `scope.requiresSchema` 优先来自 manifest，其次来自 summary。
- `status`、`message`、`requiredEnv` 来自 apply result。
- `checks`、`before`、`after` 保留现有结果。
- 默认 `actor.type` 为 `frontend-preview`。
- 默认 `eventId` 为 `preview`。
- 默认时间使用前端本地 ISO 时间，仅作为 dry-run 展示。

## 页面展示

`module-manifest-apply-result.vue` 新增 `审计预览` 操作。

点击后展示格式化 JSON 代码块，内容来自 `buildModuleManifestApplyAuditPreview`。

## 当前边界

P3.22 不接入：

- 后端接口。
- 审计表。
- 数据库写入。
- 当前登录用户读取。
- 菜单权限 SQL。

## 验收

P3.22 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
