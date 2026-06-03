# P3 Module Install Apply Boundary

更新时间：2026-06-03

## 目标

P3.10 把后台 `Manifest 预览` 的安装计划继续推进到写入门禁层。

本阶段只开放门禁检查，不执行安装 SQL，不写数据库，不创建业务 schema。

当前更新：P3.11 已按该边界开放本地受控写入 smoke；现行写入规则见 `docs/P3_MODULE_INSTALL_APPLY.md`。

## 后端接口

新增：

```text
PUT /gen/previewCode
```

该接口复用现有 `gen:previewCode` 权限面，不新增权限 SQL。

请求参数：

- `manifestPath` 或 `manifestBody`：与 `POST /gen/previewCode` 一致。
- `tenantId`：安装计划目标租户，不传或传 `0` 时使用全局租户。
- `roleId`：安装计划目标角色，不传时默认为 `1`。
- `confirmModule`：必须等于 manifest 的 `module`。
- `confirmTenantId`：必须等于本次计算出的 `tenantId`。
- `confirmRoleId`：必须等于本次计算出的 `roleId`。
- `confirmInstall`：必须为 `true`。
- `confirmSchemaRisk`：当 `manifest.requiresSchema=true` 时必须为 `true`。

写入门禁环境变量：

```bash
MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1
```

缺少环境变量或确认参数不匹配时，接口在数据库访问前失败，响应 `data` 中返回：

- manifest 摘要。
- 租户 ID、角色 ID。
- 安装计划 SQL 预览。
- 门禁检查列表。
- 阻断原因。

## P3.10 执行边界

即使环境变量和确认参数全部满足，P3.10 仍然返回失败：

```text
module install apply executor is not open in P3.10; no database access was attempted
```

原因：后台安装器需要先把写入门禁、页面确认、no-db 验证和结构化结果稳定下来，再进入本地受控写入 smoke。

## 管理端入口

`Manifest 预览` 弹窗新增：

- 确认模块输入。
- 安装写入确认。
- Schema 风险确认。
- `写入门禁` 按钮。
- 门禁检查结果表格。

页面只调用门禁接口，不执行数据库写入。

## 不在 P3.10 做

- 不执行安装 SQL。
- 不执行卸载 SQL。
- 不写 `ma_permission`、`ma_menu`、`ma_menu_permission`、`ma_role_permission`。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不新增权限 SQL。

## 验收

P3.10 需要通过：

```bash
scripts/check-module-install-apply-boundary.sh
scripts/check-module-tools-no-db.sh
cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestPreviewModuleManifest|TestModuleManifestInstallApplyGate' -count=1
cd admin && npm run type-check
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
