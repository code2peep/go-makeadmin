# P3 Module Apply Audit Preview Summary

更新时间：2026-06-03

## 目标

P3.23 将模块 manifest apply 审计事件前端预览调整为摘要优先。

本阶段只优化前端展示，不改变后端，不新增接口，不写数据库。

## 前端摘要

`admin/src/api/tools/code.ts` 新增：

```ts
buildModuleManifestApplyAuditPreviewSummary(event)
```

输出：

- 操作类型。
- 模块名。
- 执行状态。
- 路由名。
- 权限数量。
- 检查项数量。
- 执行前快照总数。
- 执行后快照总数。
- 数据库范围。
- 操作人类型。

## 页面行为

`module-manifest-apply-result.vue` 的审计预览现在分两层：

- 点击 `审计预览` 后先展示摘要。
- 点击 `JSON` 后才展开完整审计事件 JSON。
- 关闭审计预览时会同步收起 JSON。

## 当前边界

P3.23 不接入：

- 后端接口。
- 审计表。
- 数据库写入。
- 当前登录用户读取。
- 菜单权限 SQL。

## 验收

P3.23 需要通过：

```bash
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
