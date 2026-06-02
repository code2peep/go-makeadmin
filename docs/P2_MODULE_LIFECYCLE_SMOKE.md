# P2 Module Lifecycle Smoke

更新时间：2026-06-02

## 目标

P2.22 把模块安装、卸载和残留检查串成一个本地一次性 smoke 脚本，减少后续手写长 SQL 和命令的成本。

## 命令

```bash
MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 \
scripts/check-module-lifecycle-smoke.sh
```

默认数据库连接参数：

- host：`127.0.0.1`
- port：`3306`
- user：`root`
- database：`go_makeadmin`

可通过 `MYSQL_HOST`、`MYSQL_PORT`、`MYSQL_USER`、`MYSQL_DATABASE`、`MYSQL_PASSWORD` 环境变量覆盖。

## 安全门禁

脚本必须显式设置：

```bash
MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1
```

缺少环境变量时，脚本会在数据库访问前失败。

## 执行流程

脚本执行以下步骤：

1. 确认 demo article 注册行和授权残留为 0。
2. 调用 `module-install-plan.py --apply` 安装 demo article。
3. 校验安装后得到 5 条权限、1 条菜单、1 条菜单权限关联、5 条角色授权。
4. 调用 `module-uninstall-plan.py --apply` 卸载 demo article。
5. 校验卸载后残留为 0。
6. 第二次调用 `module-uninstall-plan.py --apply`，确认 no-op 幂等。

如果写入前已有 demo article 注册行，脚本会停止，不会清理已有数据。

## 不在 P2.22 做

- 不创建、修改或迁移 schema。
- 不读取或修改 `.env`。
- 不连接真实 zyai 业务库。
- 不删除文件或目录。
- 不修改 runtime 环境变量。

## 验收

P2.22 需要通过：

```bash
bash -n scripts/check-module-lifecycle-smoke.sh
scripts/check-module-lifecycle-smoke.sh
MAKEADMIN_ALLOW_MODULE_LIFECYCLE_WRITE=1 scripts/check-module-lifecycle-smoke.sh
GOCACHE=/private/tmp/go-makeadmin-gocache ./scripts/verify-no-db.sh
```
