# P3 Module Apply Audit Summary

更新时间：2026-06-03

## 目标

P3.15 为后台模块安装、卸载 apply 响应补充操作摘要，并规划后续审计日志边界。

本阶段不创建审计表，不改数据库 schema，不新增权限 SQL。

## 响应摘要

安装和卸载 apply 响应新增 `summary`：

- `operation`：`install` 或 `uninstall`。
- `module`：模块名。
- `entity`：实体名。
- `table`：来源表。
- `routeName`：菜单路由名。
- `permissionCodes`：manifest 声明的权限编码。
- `requiresSchema`：是否需要业务 schema。
- `databaseScope`：当前限定为 `local go_makeadmin only`。
- `runtimeHint`：运行时开关提示。

该摘要在成功和失败响应中都会返回，供管理端展示操作范围。

## 管理端展示

`Manifest 预览` 弹窗的安装结果和卸载结果新增：

- 操作类型。
- 路由名。
- 权限编码标签。

## 后续审计边界

未来进入审计落地时，建议单独设计审计模型，至少包含：

- 操作类型。
- manifest 来源。
- 模块名。
- 租户 ID。
- 角色 ID。
- 权限编码列表。
- 执行前快照。
- 执行后快照。
- 执行状态。
- 操作人。
- 请求时间。

审计表属于数据库 schema 变更，必须单独确认后再做。

## 不在 P3.15 做

- 不创建审计表。
- 不修改数据库 schema。
- 不执行数据库迁移。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.15 需要通过：

```bash
cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestInstallApplyGate|TestModuleManifestUninstallApplyGate' -count=1
cd admin && npm run type-check
scripts/check-module-tools-no-db.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
