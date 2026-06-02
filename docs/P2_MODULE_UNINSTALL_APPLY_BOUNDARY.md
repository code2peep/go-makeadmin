# P2 Module Uninstall Apply Boundary

更新时间：2026-06-02

## 目标

P2.20 定义模块卸载 apply/write 执行器的边界，为后续把卸载 dry-run 计划升级为本地受控删除做准备。

本阶段只写边界规范，不开放新的删除命令。

当前更新：P2.21 已按该边界开放本地受控删除模式；现行删除规则见 `docs/P2_MODULE_UNINSTALL_APPLY.md`。

## 删除门禁

未来 `module-uninstall-plan.py --apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_UNINSTALL_WRITE=1`。
- 显式传入 `--confirm-module <module>`。
- `--confirm-module` 必须等于 manifest 中的 `module`。
- 显式传入 `--confirm-delete`。

缺少环境变量或确认参数时，脚本必须在数据库访问前失败。

## 数据库范围

默认只允许本地开发库：

- host：`127.0.0.1`
- port：`3306`
- user：`root`
- database：`go_makeadmin`

可通过命令参数或 `MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_USER`、`MYSQL_DATABASE`、`MYSQL_PASSWORD` 环境变量覆盖。

执行器不得读取 `.env` 猜测数据库密码，不得连接真实 zyai 业务库。

## 删除前快照

未来 apply 删除前必须查询并输出快照计数：

- `ma_permission` 中 manifest 权限 code 命中数量。
- `ma_menu` 中 manifest 菜单 routeName 命中数量。
- `ma_menu_permission` 中引用模块菜单或模块权限的数量。
- `ma_role_permission` 中引用模块权限的数量。

如果四类计数均为 0，脚本可以报告 no-op，不执行删除。

## 删除顺序

未来 apply 必须在单事务内完成：

1. 校验 manifest。
2. 输出删除前快照计数。
3. 删除 `ma_role_permission` 中引用模块权限的授权。
4. 删除 `ma_menu_permission` 中引用模块菜单或模块权限的关联。
5. 删除 `ma_menu` 中 manifest 菜单 routeName 对应节点。
6. 删除 `ma_permission` 中 manifest 权限 code 对应权限。

事务内任一步失败，全部回滚。

## 删除边界

- 只按 manifest 声明的权限 code 和菜单 routeName 删除。
- 不删除前端文件。
- 不删除后端代码。
- 不关闭 runtime 环境变量。
- 不删除 schema。
- 不删除 `ma_codegen_table` 或 `ma_codegen_column`。

## 本地 smoke 要求

未来开放 apply 后，必须对本地 `go_makeadmin` 执行一次 demo article smoke：

- 写入前确认 demo article 注册行不存在。
- 先用 `module-install-plan.py --apply` 安装 demo article。
- 再用 `module-uninstall-plan.py --apply` 卸载 demo article。
- 卸载后残留计数为 0。
- 再执行第二次 uninstall apply，确认 no-op 幂等。

如果写入前已有 demo article 注册行，smoke 必须停止，不得删除已有行。

## 不在 P2.20 做

- 不开放 `--apply`。
- 不执行数据库删除。
- 不删除文件或目录。
- 不创建、修改或迁移 schema。
- 不修改 `.env`、密钥、CI/CD 或生产配置。
- 不启动服务或部署。

## 验收

P2.20 需要通过：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
