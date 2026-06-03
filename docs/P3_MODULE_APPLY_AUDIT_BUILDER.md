# P3 Module Apply Audit Builder

更新时间：2026-06-03

## 目标

P3.21 为模块 manifest 安装、卸载 apply 建立审计事件构造器 dry-run。

本阶段只做纯函数和单测，不接入接口，不创建审计表，不写数据库。

## 构造器

新增：

```text
server/generator/service/gen/module_manifest_apply_audit.go
```

包含：

- `buildModuleManifestInstallAuditEvent`
- `buildModuleManifestUninstallAuditEvent`
- `buildModuleManifestApplyAuditEvent`

构造器输入：

- apply result。
- event ID。
- actor。
- requestedAt。
- completedAt。

构造器输出：

- `resp.ModuleManifestApplyAuditEventResp`。

## 映射规则

- `source` 来自 apply result。
- `manifest` 来自 apply result。
- `summary` 来自 apply result。
- `scope.tenantId` 和 `scope.roleId` 来自 install/uninstall plan。
- `scope.databaseScope` 来自 apply summary。
- `scope.requiresSchema` 来自 manifest summary。
- `status`、`message`、`requiredEnv` 来自 apply result。
- `checks`、`before`、`after` 原样保留。
- `actor`、`requestedAt`、`completedAt` 由调用方显式传入。
- 构造器不生成 event ID，不读取数据库，不读取当前用户。

## 单测

新增：

```text
server/generator/service/gen/module_manifest_apply_audit_test.go
```

覆盖：

- 安装 apply result 到审计事件。
- 卸载 apply result 到审计事件。
- 执行范围、操作类型、操作人、快照和门禁检查映射。

## 当前边界

P3.21 不接入：

- `PUT /gen/previewCode`
- `DELETE /gen/previewCode`
- 审计表
- 数据库写入
- 菜单权限 SQL
- 前端展示

## 后续落地建议

下一步可以做前端 dry-run 预览，把当前页面中的 apply result 组合成审计事件预览，但仍不接入写库。

## 验收

P3.21 需要通过：

```bash
cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestBuildModuleManifest.*AuditEvent' -count=1
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
