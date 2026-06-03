# P3 Module Install Apply

更新时间：2026-06-03

## 目标

P3.11 在 P3.10 写入门禁基础上，开放后台模块安装的本地受控写入 smoke。

本阶段只针对本地 `go_makeadmin` 开发库，写入 manifest 中声明的：

- `ma_permission`
- `ma_menu`
- `ma_menu_permission`
- `ma_role_permission`

不创建、修改或迁移业务 schema。

## 后端接口

接口仍为：

```text
PUT /gen/previewCode
```

继续复用 `gen:previewCode` 权限面，不新增权限 SQL。

## 写入门禁

接口写入必须满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_INSTALL_APPLY=1`。
- `confirmModule` 等于 manifest 的 `module`。
- `confirmTenantId` 等于本次安装计划的 `tenantId`。
- `confirmRoleId` 等于本次安装计划的 `roleId`。
- `confirmInstall=true`。
- `manifest.requiresSchema=true` 时必须 `confirmSchemaRisk=true`。
- 服务配置的数据库 DSN 必须指向本地 `go_makeadmin`。

缺少环境变量、确认参数不匹配、数据库目标不是本地 `go_makeadmin` 或服务未配置数据库时，接口会在数据库访问前失败，并返回 `no database access was attempted`。

## 写入规则

- 所有写入在单事务内完成。
- 权限按 `ma_permission.code` 幂等。
- 菜单按 `ma_menu.route_name + delete_time=0` 幂等。
- 菜单权限按 `menu_id + permission_id` 幂等。
- 角色授权按 `tenant_id + role_id + permission_id` 幂等。
- 已存在行不覆盖、不更新。
- 如果目标角色不存在，角色授权跳过，菜单和权限仍可安装。
- runtime 开关只返回提示，不修改 `.env` 或系统环境变量。

## Smoke

新增：

```bash
MAKEADMIN_ALLOW_MODULE_INSTALL_SMOKE_WRITE=1 \
scripts/check-module-install-apply-smoke.sh
```

smoke 流程：

1. 确认 demo article 注册行和角色授权残留为 0。
2. 调用后台服务执行第一次安装。
3. 验证写入 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
4. 调用后台服务执行第二次安装。
5. 验证计数不变，确认幂等。
6. 清理本次 demo article 写入的数据。
7. 验证清理后残留为 0。

如果写入前已有 demo article 相关行，smoke 会停止，不会清理已有数据。

## 不在 P3.11 做

- 不执行卸载接口。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不连接真实业务项目数据库。
- 不新增权限 SQL。

## 验收

P3.11 需要通过：

```bash
scripts/check-module-install-apply-boundary.sh
scripts/check-module-tools-no-db.sh
MAKEADMIN_ALLOW_MODULE_INSTALL_SMOKE_WRITE=1 scripts/check-module-install-apply-smoke.sh
cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestInstallApplyGate|TestModuleManifestInstallApplyLocalSmoke' -count=1
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
