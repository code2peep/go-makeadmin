# P2 Module Install Apply Boundary

更新时间：2026-06-02

## 目标

P2.17 定义模块安装 apply/write 执行器的边界，为后续把 dry-run 安装计划升级为本地受控写入做准备。

本阶段只写边界规范，不开放新的写入命令。

当前更新：P2.18 已按该边界开放本地受控写入模式；现行写入规则见 `docs/P2_MODULE_INSTALL_APPLY.md`。

## 写入门禁

未来 `module-install-plan.py --apply` 必须同时满足：

- 环境变量 `MAKEADMIN_ALLOW_MODULE_INSTALL_WRITE=1`。
- 显式传入 `--confirm-module <module>`。
- `--confirm-module` 必须等于 manifest 中的 `module`。
- 如果传入 `--role-id`，还必须显式传入 `--confirm-role-id <id>`。
- `--confirm-role-id` 必须等于 `--role-id`。

缺少环境变量或确认参数时，脚本必须在数据库访问前失败。

## 数据库范围

默认只允许本地开发库：

- host：`127.0.0.1`
- port：`3306`
- user：`root`
- database：`go_makeadmin`

可通过命令参数或 `MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_USER`、`MYSQL_DATABASE`、`MYSQL_PASSWORD` 环境变量覆盖。

执行器不得读取 `.env` 猜测数据库密码，不得连接真实 zyai 业务库。

## 写入顺序

未来 apply 必须在单事务内完成：

1. 校验 manifest。
2. 校验 `requiresSchema=false`；如果为 `true`，本阶段直接失败，不自动建表。
3. 写入缺失 `ma_permission`。
4. 写入缺失 `ma_menu`。
5. 写入缺失 `ma_menu_permission`。
6. 如果传入 `--role-id`，写入缺失 `ma_role_permission`。

事务内任一步失败，全部回滚。

## 幂等规则

- 权限按 `ma_permission.code` 幂等。
- 菜单按 `ma_menu.route_name + delete_time=0` 幂等。
- 菜单权限按 `menu_id + permission_id` 幂等。
- 角色授权按 `tenant_id + role_id + permission_id` 幂等。
- 已存在行不覆盖、不更新。

## runtime 边界

runtime 开关不通过数据库写入改变。

当前 demo runtime 仍通过环境变量控制：

```bash
MAKEADMIN_ENABLE_DEMO_MODULE=1
```

安装执行器只能输出 runtime 提示，不能修改 `.env`、系统环境变量或生产配置。

## 本地 smoke 要求

未来开放 apply 后，必须对本地 `go_makeadmin` 执行一次 demo article smoke：

- 写入前确认 demo article 注册行不存在。
- 第一次 apply 插入权限、菜单、菜单权限和可选角色授权。
- 第二次 apply 计数不变，确认幂等。
- 清理 demo article 注册行。
- 清理后残留计数为 0。

如果写入前已有 demo article 注册行，smoke 必须停止，不得删除已有行。

## 不在 P2.17 做

- 不开放 `--apply`。
- 不执行数据库写入。
- 不创建、修改或迁移 schema。
- 不修改 `.env`、密钥、CI/CD 或生产配置。
- 不启动服务或部署。

## 验收

P2.17 需要通过：

```bash
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
