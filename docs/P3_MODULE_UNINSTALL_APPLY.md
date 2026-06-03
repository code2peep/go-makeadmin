# P3 Module Uninstall Apply

更新时间：2026-06-03

## 目标

P3.13 在 P3.12 卸载门禁基础上，开放后台模块卸载的本地受控删除 smoke。

本阶段只针对本地 `go_makeadmin` 开发库，删除 manifest 中声明的：

- `ma_role_permission`
- `ma_menu_permission`
- `ma_menu`
- `ma_permission`

不删除业务表，不创建、修改或迁移业务 schema。

## 后端接口

接口仍为：

```text
DELETE /gen/previewCode
```

继续复用 `gen:previewCode` 权限面，不新增权限 SQL。

## 写入门禁

接口删除必须满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_APPLY=1`。
- `confirmModule` 等于 manifest 的 `module`。
- `confirmDelete=true`。
- 服务配置的数据库 DSN 必须指向本地 `go_makeadmin`。

缺少环境变量、确认参数不匹配、数据库目标不是本地 `go_makeadmin` 或服务未配置数据库时，接口会在数据库访问前失败，并返回 `no database access was attempted`。

## 删除规则

- 所有删除在单事务内完成。
- 先删除角色授权，再删除菜单权限，再删除菜单，最后删除权限。
- 权限按 manifest 声明的 permission codes 删除。
- 菜单按 manifest 声明的 `menu.routeName` 删除。
- 菜单权限按 `menu.routeName` 或 permission codes 删除。
- 角色授权按 permission codes 删除。
- 第二次执行应为 no-op，最终残留仍为 0。
- 不删除业务表、不修改 runtime 配置。

## Smoke

新增：

```bash
MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE=1 \
scripts/check-module-uninstall-apply-smoke.sh
```

smoke 流程：

1. 确认 demo article 注册行和角色授权残留为 0。
2. 调用后台安装接口种入 demo article。
3. 调用后台卸载接口执行第一次删除。
4. 验证权限、菜单、菜单权限和角色授权残留为 0。
5. 第二次调用后台卸载接口。
6. 验证 no-op 幂等，残留仍为 0。

如果写入前已有 demo article 相关行，smoke 会停止，不会清理已有数据。

## 不在 P3.13 做

- 不删除业务表。
- 不创建、修改或迁移业务 schema。
- 不读取或修改 `.env`、密钥、CI/CD 或生产配置。
- 不连接真实业务项目数据库。
- 不新增权限 SQL。

## 验收

P3.13 需要通过：

```bash
scripts/check-module-uninstall-apply-boundary.sh
scripts/check-module-tools-no-db.sh
MAKEADMIN_ALLOW_MODULE_UNINSTALL_SMOKE_WRITE=1 scripts/check-module-uninstall-apply-smoke.sh
cd server && GOCACHE=/private/tmp/go-makeadmin-gocache go test ./generator/service/gen -run 'TestModuleManifestUninstallApplyGate|TestModuleManifestUninstallApplyLocalSmoke' -count=1
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
