# P3 Module Apply Audit DTO

更新时间：2026-06-03

## 目标

P3.20 为模块 manifest 安装、卸载 apply 审计事件建立 DTO 草图。

本阶段只定义 Go/TypeScript 数据形状和文档边界，不创建审计表，不接入写库，不修改接口响应。

## Go DTO

新增：

```text
server/generator/schemas/resp/module_manifest_audit.go
```

包含：

- `ModuleManifestApplyAuditActorResp`
- `ModuleManifestApplyAuditScopeResp`
- `ModuleManifestApplyAuditEventResp`

审计事件草图覆盖：

- 事件 ID。
- 操作类型。
- manifest 来源。
- manifest 摘要。
- 操作摘要。
- 租户、角色、数据库范围和 schema 风险。
- 执行状态和说明。
- 环境变量门禁。
- 检查结果。
- 执行前后快照。
- 操作人。
- 请求时间和完成时间。

## TypeScript DTO

`admin/src/api/tools/code.ts` 新增：

- `ModuleManifestApplyAuditActorResult`
- `ModuleManifestApplyAuditScopeResult`
- `ModuleManifestApplyAuditEventResult`

字段名按后端 JSON 命名，供后续管理端审计页面、结果详情或接口草图复用。

## 当前边界

P3.20 的 DTO 不接入：

- `PUT /gen/previewCode`
- `DELETE /gen/previewCode`
- `ModuleManifestApplyResult`
- 数据库模型
- 数据库迁移
- 审计写入服务
- 菜单权限 SQL

## 后续落地建议

如果进入真实审计落地，建议单独确认：

- 审计表 schema。
- 审计事件 ID 生成规则。
- 操作人来源。
- 是否记录完整 SQL 预览。
- 是否记录完整请求体。
- 是否支持审计列表、详情和导出。
- 审计保留时长。

审计表属于数据库 schema 变更，必须单独确认后再做。

## 验收

P3.20 需要通过：

```bash
gofmt -w server/generator/schemas/resp/module_manifest_audit.go
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
