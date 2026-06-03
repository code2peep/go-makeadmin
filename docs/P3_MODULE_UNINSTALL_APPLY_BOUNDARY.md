# P3 Module Uninstall Apply Boundary

更新时间：2026-06-03

## 目标

P3.12 为后台模块卸载建立写入门禁和确认参数。

本阶段只开放门禁检查，不执行卸载 SQL，不删除数据库行，不创建或迁移 schema。

## 后端接口

新增：

```text
DELETE /gen/previewCode
```

该接口复用现有 `gen:previewCode` 权限面，不新增权限 SQL。

请求参数：

- `manifestPath` 或 `manifestBody`：与 `POST /gen/previewCode` 一致。
- `confirmModule`：必须等于 manifest 的 `module`。
- `confirmDelete`：必须为 `true`。

写入门禁环境变量：

```bash
MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1
```

缺少环境变量或确认参数不匹配时，接口在数据库访问前失败，响应 `data` 中返回：

- manifest 摘要。
- 卸载 SQL 预览。
- 门禁检查列表。
- 阻断原因。

## P3.12 执行边界

即使环境变量和确认参数全部满足，P3.12 仍然返回失败：

```text
module uninstall apply executor is not open in P3.12; no database access was attempted
```

原因：后台卸载器需要先固定门禁和结构化结果，再进入本地受控删除 smoke。

## 不在 P3.12 做

- 不执行卸载 SQL。
- 不删除 `ma_permission`、`ma_menu`、`ma_menu_permission`、`ma_role_permission`。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.12 需要通过：

```bash
scripts/check-module-uninstall-apply-boundary.sh
scripts/check-module-tools-no-db.sh
cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestUninstallApplyGate' -count=1
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
